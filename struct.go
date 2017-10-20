// cache struct info to improve encoding/decoding efficiency.

package binary

import (
	"fmt"
	"reflect"
)

// RegStruct regist struct info to improve encoding/decoding efficiency.
// Regist by a nil pointer is aviable.
func RegStruct(data interface{}) error {
	return _structInfoMgr.regist(reflect.TypeOf(data))
}

var _structInfoMgr structInfoMgr

func init() {
	_structInfoMgr.init()
}

//var intNameToType = map[string]reflect.Kind{
//	"uint8":   reflect.Uint8,
//	"uint16":  reflect.Uint16,
//	"uint32":  reflect.Uint32,
//	"uint64":  reflect.Uint64,
//	"int8":    reflect.Int8,
//	"int16":   reflect.Int16,
//	"int32":   reflect.Int32,
//	"int64":   reflect.Int64,
//	"int":     reflect.Int,
//	"uint":    reflect.Uint,
//	"varint":  reflect.Int,
//	"uvarint": reflect.Uint,
//}

//func getIntKind(kind string) reflect.Kind {
//	if k, ok := intNameToType[kind]; ok {
//		return k
//	}
//	//panic("binary: unsupported int kind " + kind)
//	return reflect.Invalid
//}

type structInfoMgr struct {
	reg map[string]*structInfo
}

func (this *structInfoMgr) init() {
	this.reg = make(map[string]*structInfo)
}
func (this *structInfoMgr) regist(t reflect.Type) error {
	if _t, _, err := this.deepStructType(t, true); err == nil {
		if this.query(_t) == nil {
			p := &structInfo{}
			if p.parse(_t) {
				this.reg[p.identify] = p
			}
		} else {
			return fmt.Errorf("binary: regist duplicate type %s", _t.String())
		}
	} else {
		return err
	}
	return nil
}
func (this *structInfoMgr) query(t reflect.Type) *structInfo {
	if _t, _ok, _ := this.deepStructType(t, false); _ok {
		if p, ok := this.reg[_t.String()]; ok {
			return p
		}
	}
	return nil
}

func (this *structInfoMgr) deepStructType(t reflect.Type, needErr bool) (reflect.Type, bool, error) {
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

func (this *structInfo) encode(encoder *Encoder, v reflect.Value) error {
	assert(v.Kind() == reflect.Struct, v.Type().String())
	t := v.Type()
	for i, n := 0, v.NumField(); i < n; i++ {
		// see comment for corresponding code in decoder.value()
		if f := v.Field(i); this.fieldValid(i, t) {
			if err := encoder.value(f); err != nil {
				return err
			}
		} else {
			//do nothing
		}
	}
	return nil
}

//func (this *structInfo) encodeField(encoder *Encoder, i int, v reflect.Value) error {}
//func (this *structInfo) decodeField(decoder *Decoder, i int, v reflect.Value) error {}

func (this *structInfo) encodeNilPointer(encoder *Encoder, t reflect.Type) int {
	sum := 0
	for i, n := 0, this.fieldNum(t); i < n; i++ {
		s := encoder.nilPointer(this.fieldType(i, t))
		if s < 0 {
			return -1
		}
		sum += s
	}
	return sum
}

func (this *structInfo) decode(decoder *Decoder, v reflect.Value) error {
	assert(v.Kind() == reflect.Struct, v.Type().String())
	t := v.Type()
	for i, n := 0, v.NumField(); i < n; i++ {
		if f := v.Field(i); this.fieldValid(i, t) {
			if err := decoder.value(f); err != nil {
				return err
			}
		} else {
			//do nothing
		}
	}
	return nil
}

func (this *structInfo) decodeSkipByType(decoder *Decoder, t reflect.Type) int {
	sum := 0
	for i, n := 0, t.NumField(); i < n; i++ {
		ft := this.fieldType(i, t)
		s := decoder.skipByType(ft)
		assert(s >= 0, ft.String()) //I'm sure here cannot find unsupported type
		sum += s
	}
	return sum
}

func (this *structInfo) sizeofValue(v reflect.Value) int {
	assert(v.Kind() == reflect.Struct, v.Type().String())
	t := v.Type()
	sum := 0
	for i, n := 0, v.NumField(); i < n; i++ {
		if this.fieldValid(i, t) {
			s := sizeofValue(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
	}
	return sum
}

func (this *structInfo) sizeofEmptyPointer(t reflect.Type) int {
	sum := 0
	for i, n := 0, this.fieldNum(t); i < n; i++ {
		s := sizeofNilPointer(this.fieldType(i, t))
		if s < 0 {
			return -1
		}
		sum += s
	}
	return sum
}

//check if
func (this *structInfo) fieldValid(i int, t reflect.Type) bool {
	if this == nil {
		//Note: creating the StructField info for each field is costly
		return validField(t.Field(i)) // slow way to access field info
	} else {
		return this.field(i).valid() //fast way to access field info
	}
}

func (this *structInfo) fieldType(i int, t reflect.Type) reflect.Type {
	if this == nil {
		return t.Field(i).Type
	} else {
		return this.field(i).field.Type
	}
}
func (this *structInfo) fieldNum(t reflect.Type) int {
	if this == nil {
		return t.NumField()
	} else {
		return this.numField()
	}
}

func (this *structInfo) parse(t reflect.Type) bool {
	assert(t.Kind() == reflect.Struct, t.String())
	this.identify = t.String()
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)

		field := &fieldInfo{}
		field.field = f
		tag := f.Tag.Get("binary")
		field.ignore = !isExported(f.Name) || tag == "ignore"
		//field.encodeKind = getIntKind(tag)

		this.fields = append(this.fields, field)

		//deep regist if field is a struct
		if _t, ok, _ := _structInfoMgr.deepStructType(f.Type, false); ok {
			if err := _structInfoMgr.regist(_t); err != nil {
				//fmt.Printf("binary: internal regist duplicate type %s\n", _t.String())
			}
		}
	}
	return true
}

func (this *structInfo) field(i int) *fieldInfo {
	if nil != this && i >= 0 && i < this.numField() {
		return this.fields[i]
	}
	return nil
}

func (this *structInfo) numField() int {
	if nil != this {
		return len(this.fields)
	}
	return 0
}

//informatin of a struct field
type fieldInfo struct {
	field  reflect.StructField
	ignore bool //if this field is ignored
	//encodeKind reflect.Kind //enable encode integers as other size
}

func (this *fieldInfo) valid() bool {
	return !this.ignore
}

func queryStruct(t reflect.Type) *structInfo {
	return _structInfoMgr.query(t)
}
