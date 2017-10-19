//cache struct info to improve encoding/decoding efficiency

package binary

/*

import (
	"fmt"
	"reflect"
)

// RegsterStruct regist struct info to improve encoding/decoding efficiency
func RegistStruct(data interface{}) error {
	return _structInfoMgr.regist(reflect.TypeOf(data))
}

var _structInfoMgr structInfoMgr

func init() {
	_structInfoMgr.init()
}

var intNameToType = map[string]reflect.Kind{
	"uint8":   reflect.Uint8,
	"uint16":  reflect.Uint16,
	"uint32":  reflect.Uint32,
	"uint64":  reflect.Uint64,
	"int8":    reflect.Int8,
	"int16":   reflect.Int16,
	"int32":   reflect.Int32,
	"int64":   reflect.Int64,
	"int":     reflect.Int,
	"uint":    reflect.Uint,
	"varint":  reflect.Int,
	"uvarint": reflect.Uint,
}

func getIntKind(kind string) reflect.Kind {
	if k, ok := intNameToType[kind]; ok {
		return k
	}
	panic("binary: unsupported int kind " + kind)
	//return reflect.Invalid
}

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
	return nil
}
func (this *structInfo) decode(decoder *Decoder, v reflect.Value) error {
	return nil
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
		field.encodeKind = getIntKind(tag)

		this.fields = append(this.fields, field)

		//deep register if field is a struct
		if _t, ok, _ := _structInfoMgr.deepStructType(f.Type, false); ok {
			if err := _structInfoMgr.regist(_t); err != nil {
				//fmt.Printf("binary: internal regist duplicate type %s\n", _t.String())
			}
		}
	}
	return true
}

func (this *structInfo) field(i int) *fieldInfo {
	if i >= 0 && i < this.numField() {
		return this.fields[i]
	}
	return nil
}

func (this *structInfo) numField() int {
	return len(this.fields)
}

//informatin of a struct field
type fieldInfo struct {
	field      reflect.StructField
	ignore     bool         //if this field is ignored
	encodeKind reflect.Kind //what this field will be encoded
}

func (this *fieldInfo) valid() bool {
	return !this.ignore
}

func queryStruct(t reflect.Type) *structInfo {
	return _structInfoMgr.query(t)
}

//*/
