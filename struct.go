// cache struct info to improve encoding/decoding efficiency.

package binary

import (
	"fmt"
	"reflect"
)

// RegStruct regist struct info to improve encoding/decoding efficiency.
// Regist by a nil pointer is aviable.
// RegStruct((*someStruct)(nil)) is recommended usage.
func RegStruct(data interface{}) error {
	return _structInfoMgr.regist(reflect.TypeOf(data))
}

var _structInfoMgr structInfoMgr

func init() {
	_structInfoMgr.init()
}

type structInfoMgr struct {
	reg map[string]*structInfo
}

func (mgr *structInfoMgr) init() {
	mgr.reg = make(map[string]*structInfo)
}
func (mgr *structInfoMgr) regist(t reflect.Type) error {
	if _t, _, err := mgr.deepStructType(t, true); err == nil {
		if mgr.query(_t) == nil {
			p := &structInfo{}
			if p.parse(_t) {
				mgr.reg[p.identify] = p
			}
		} else {
			return fmt.Errorf("binary: regist duplicate type %s", _t.String())
		}
	} else {
		return err
	}
	return nil
}

func (mgr *structInfoMgr) query(t reflect.Type) *structInfo {
	if _t, _ok, _ := mgr.deepStructType(t, false); _ok {
		if p, ok := mgr.reg[_t.String()]; ok {
			return p
		}
	}
	return nil
}

func (mgr *structInfoMgr) deepStructType(t reflect.Type, needErr bool) (reflect.Type, bool, error) {
	_t := t
	for _t.Kind() == reflect.Ptr {
		_t = _t.Elem()
	}
	if _t.Kind() != reflect.Struct {
		if needErr {
			return _t, false, fmt.Errorf("binary: only struct is aviable for regist, but got %s", t.String())
		} else {
			return _t, false, nil
		}
	}
	return _t, true, nil
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
		if f := v.Field(i); finfo.valid(i, t) {
			if err := encoder.value(f, finfo.packed()); err != nil {
				return err
			}
		} else {
			//do nothing
		}
	}
	return nil
}

func (info *structInfo) decode(decoder *Decoder, v reflect.Value) error {
	t := v.Type()
	//assert(t.Kind() == reflect.Struct, t.String())
	for i, n := 0, v.NumField(); i < n; i++ {
		finfo := info.field(i)
		if f := v.Field(i); finfo.valid(i, t) {
			if err := decoder.value(f, false, finfo.packed()); err != nil {
				return err
			}
		} else {
			//do nothing
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
		s := decoder.skipByType(ft, f.packed())
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

		if finfo := info.field(i); finfo.valid(i, t) {
			if s := bitsOfValue(v.Field(i), false, finfo.packed()); s >= 0 {
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
	return info.field(i).valid(i, t)
}

func (info *structInfo) fieldNum(t reflect.Type) int {
	if info == nil {
		return t.NumField()
	} else {
		return info.numField()
	}
}

func (info *structInfo) parse(t reflect.Type) bool {
	//assert(t.Kind() == reflect.Struct, t.String())
	info.identify = t.String()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)

		field := &fieldInfo{}
		field.field = f
		tag := f.Tag.Get("binary")
		field.ignore = !isExported(f.Name) || tag == "ignore"
		field.packed_ = tag == "packed"

		info.fields = append(info.fields, field)

		//deep regist if field is a struct
		if _t, ok, _ := _structInfoMgr.deepStructType(f.Type, false); ok {
			if err := _structInfoMgr.regist(_t); err != nil {
				//fmt.Printf("binary: internal regist duplicate type %s\n", _t.String())
			}
		}
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
	field   reflect.StructField
	ignore  bool //if this field is ignored
	packed_ bool //if this ints field encode as varint/uvarint
}

func (field *fieldInfo) Type(i int, t reflect.Type) reflect.Type {
	if field != nil {
		return field.field.Type
	} else {
		return t.Field(i).Type
	}
}

func (field *fieldInfo) valid(i int, t reflect.Type) bool {
	if field != nil {
		return !field.ignore
	} else {
		// NOTE:
		// creating the StructField info for each field is costly
		// use RegStruct((*someStruct)(nil)) to aboid this path
		return validField(t.Field(i)) // slow way to access field info
	}
}

func (field *fieldInfo) packed() bool {
	return field != nil && field.packed_
}

func queryStruct(t reflect.Type) *structInfo {
	return _structInfoMgr.query(t)
}
