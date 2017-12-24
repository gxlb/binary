// cache struct info to improve encoding/decoding efficiency.
// regist serializer type to improve type checking efficiency.

package binary

import (
	"fmt"
	"reflect"
)

// RegsterType regist type info to improve encoding/decoding efficiency.
// Only BinarySerializer or struct is regable.
// Regist by a nil pointer is aviable.
// RegsterType((*SomeType)(nil)) is recommended usage.
func RegsterType(x interface{}) error {
	return _regedTypeMgr.regist(reflect.TypeOf(x), true)
}

var (
	//tSizer        reflect.Type //BinarySizer
	//tDecoder      reflect.Type //BinaryDecoder
	tEncoder      reflect.Type //BinaryEncoder
	tSerializer   reflect.Type //BinarySerializer
	_regedTypeMgr regedTypeMgr //reged type manager
)

func init() {
	//var sizer BinarySizer
	//var decoder BinaryDecoder
	//tSizer = reflect.TypeOf(&sizer).Elem()
	//tDecoder = reflect.TypeOf(&decoder).Elem()
	var encoder BinaryEncoder
	var serializer BinarySerializer
	tEncoder = reflect.TypeOf(&encoder).Elem()
	tSerializer = reflect.TypeOf(&serializer).Elem()

	_regedTypeMgr.init()
}

type regedTypeMgr struct {
	regedStruct     map[reflect.Type]*structInfo
	regedSerializer map[reflect.Type]bool
}

func (mgr *regedTypeMgr) init() {
	mgr.regedStruct = make(map[reflect.Type]*structInfo)
	mgr.regedSerializer = make(map[reflect.Type]bool)
}

func (mgr *regedTypeMgr) regist(t reflect.Type, needError bool) (err error) {
	_t, isSerializer, ok, _err := mgr.deepRegableType(t, needError)
	if err = _err; ok {
		if _t.Kind() == reflect.Struct {
			err = mgr.regstruct(_t, needError)
		}
		if isSerializer {
			err = mgr.regserializer(_t, needError)
		}
	}
	return
}

func (mgr *regedTypeMgr) regstruct(t reflect.Type, needError bool) error {
	if mgr.queryStruct(t) == nil {
		p := &structInfo{}
		if p.parse(mgr, t) {
			mgr.regedStruct[t] = p
		}
		needError = false
	}
	return typeError("binary: regist duplicate type %s", t, needError)
}
func (mgr *regedTypeMgr) regserializer(t reflect.Type, needError bool) error {
	if !mgr.querySerializer(t) {
		mgr.regedSerializer[t] = true
		needError = false

		//reg sub data type for data-set
		switch t.Kind() {
		case reflect.Struct: //struct has reged by regstruct
		case reflect.Map:
			mgr.regist(t.Key(), false)
			mgr.regist(t.Elem(), false)
		case reflect.Slice, reflect.Array:
			mgr.regist(t.Elem(), false)
		}
	}

	return typeError("binary: regist duplicate BinarySerializer %s", t, needError)
}

func (mgr *regedTypeMgr) querySerializer(t reflect.Type) bool {
	_, ok := mgr.regedSerializer[t]
	return ok
}

func (mgr *regedTypeMgr) queryStruct(t reflect.Type) *structInfo {
	if p, ok := mgr.regedStruct[t]; ok {
		return p
	}
	return nil
}

func typeError(fmt_ string, t reflect.Type, needErr bool) error {
	if needErr {
		return fmt.Errorf(fmt_, t.String())
	}
	return nil
}

func (mgr *regedTypeMgr) deepRegableType(t reflect.Type, needErr bool) (deept reflect.Type, isSerializer, ok bool, err error) {
	_t := t
	for _t.Kind() == reflect.Ptr {
		_t = _t.Elem()
	}

	if _pt := reflect.PtrTo(_t); _t.Implements(tEncoder) || _pt.Implements(tEncoder) {
		if !_pt.Implements(tSerializer) {
			return t, false, false, typeError("binary: unexpected BinaryEncoder, expect implements BinarySerializer, type %s", t, needErr)
		}
		isSerializer = true
	}

	if isSerializer || _t.Kind() == reflect.Struct {
		return _t, isSerializer, true, nil
	}

	return t, false, false, typeError("binary: expect Regist by BinarySerializer or struct, type %s", t, needErr)
}

//informatin of a struct
type structInfo struct {
	identify string //reflect.Type.String()
	fields   []*fieldInfo
}

func (info *structInfo) encode(encoder *Encoder, v reflect.Value) error {
	//assert(v.Kind() == reflect.Struct, v.Type().String())
	t := v.Type()
	for i, n := 0, v.NumField(); i < n; i++ {
		// see comment for corresponding code in decoder.value()
		finfo := info.field(i)
		if f := v.Field(i); finfo.isValid(i, t) {
			if err := encoder.value(f, finfo.isPacked()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (info *structInfo) decode(decoder *Decoder, v reflect.Value) error {
	t := v.Type()
	//assert(t.Kind() == reflect.Struct, t.String())
	for i, n := 0, v.NumField(); i < n; i++ {
		finfo := info.field(i)
		if f := v.Field(i); finfo.isValid(i, t) {
			if err := decoder.value(f, false, finfo.isPacked()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (info *structInfo) decodeSkipByType(decoder *Decoder, t reflect.Type, packed bool) int {
	//assert(t.Kind() == reflect.Struct, t.String())
	sum := 0
	for i, n := 0, t.NumField(); i < n; i++ {
		f := info.field(i)
		ft := f.Type(i, t)
		s := decoder.skipByType(ft, f.isPacked())
		assert(s >= 0, "skip struct field fail:"+ft.String()) //I'm sure here cannot find unsupported type
		sum += s
	}
	return sum
}

func (info *structInfo) bitsOfValue(v reflect.Value) int {
	t := v.Type()
	//assert(t.Kind() == reflect.Struct,t.String())
	sum := 0
	for i, n := 0, v.NumField(); i < n; i++ {

		if finfo := info.field(i); finfo.isValid(i, t) {
			if s := bitsOfValue(v.Field(i), false, finfo.isPacked()); s >= 0 {
				sum += s
			} else {
				return -1 //invalid field type
			}
		}
	}
	return sum
}

func (info *structInfo) sizeofNilPointer(t reflect.Type) int {
	sum := 0
	for i, n := 0, info.fieldNum(t); i < n; i++ {
		if info.fieldValid(i, t) {
			if s := sizeofNilPointer(info.field(i).Type(i, t)); s >= 0 {
				sum += s
			} else {
				return -1 //invalid field type
			}
		}
	}
	return sum
}

//check if field i of t valid for encoding/decoding
func (info *structInfo) fieldValid(i int, t reflect.Type) bool {
	return info.field(i).isValid(i, t)
}

func (info *structInfo) fieldNum(t reflect.Type) int {
	if info == nil {
		return t.NumField()
	}

	return info.numField()
}

func (info *structInfo) parse(mgr *regedTypeMgr, t reflect.Type) bool {
	//assert(t.Kind() == reflect.Struct, t.String())
	info.identify = t.String()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)

		field := &fieldInfo{}
		field.field = f
		tag := f.Tag.Get("binary")
		field.ignore = !isExported(f.Name) || tag == "ignore"
		field.packed = tag == "packed"

		info.fields = append(info.fields, field)

		//deep regist field type
		mgr.regist(f.Type, false)
	}
	return true
}

func (info *structInfo) field(i int) *fieldInfo {
	if nil != info && i >= 0 && i < info.numField() {
		return info.fields[i]
	}
	return nil
}

func (info *structInfo) numField() int {
	if nil != info {
		return len(info.fields)
	}
	return 0
}

//informatin of a struct field
type fieldInfo struct {
	field  reflect.StructField
	ignore bool //if this field is ignored
	packed bool //if this ints field encode as varint/uvarint
}

func (field *fieldInfo) Type(i int, t reflect.Type) reflect.Type {
	if field != nil {
		return field.field.Type
	}

	return t.Field(i).Type
}

func (field *fieldInfo) isValid(i int, t reflect.Type) bool {
	if field != nil {
		return !field.ignore
	}

	// NOTE:
	// creating the StructField info for each field is costly
	// use RegStruct((*someStruct)(nil)) to aboid this path
	return validField(t.Field(i)) // slow way to access field info
}

func (field *fieldInfo) isPacked() bool {
	return field != nil && field.packed
}

func queryStruct(t reflect.Type) *structInfo {
	return _regedTypeMgr.queryStruct(t)
}

func querySerializer(t reflect.Type) bool {
	return _regedTypeMgr.querySerializer(t)
}
