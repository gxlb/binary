package binary

import (
	"fmt"
	"reflect"
)

const (
	//disable BinarySerializer check by default
	defaultSerializer = true
)

//DefaultSerializer return default BinarySerializer check by default
func DefaultSerializer() bool {
	return defaultSerializer
}

// BinarySizer is an interface to define go data Size method.
type BinarySizer interface {
	Size() int
}

// BinaryEncoder is an interface to define go data Encode method.
// buffer is nil-able.
type BinaryEncoder interface {
	Encode(buffer []byte) ([]byte, error)
}

// BinaryDecoder is an interface to define go data Decode method.
type BinaryDecoder interface {
	Decode(buffer []byte) error
}

// BinarySerializer defines the go data Size/Encode/Decode method
type BinarySerializer interface {
	BinarySizer
	BinaryEncoder
	BinaryDecoder
}

var (
	tSizer      reflect.Type //BinarySizer
	tEncoder    reflect.Type //BinaryEncoder
	tDecoder    reflect.Type //BinaryDecoder
	tSerializer reflect.Type //BinarySerializer
)

func init() {
	var sizer BinarySizer
	tSizer = reflect.TypeOf(&sizer).Elem()
	var encoder BinaryEncoder
	tEncoder = reflect.TypeOf(&encoder).Elem()
	var decoder BinaryDecoder
	tDecoder = reflect.TypeOf(&decoder).Elem()
	var serializer BinarySerializer
	tSerializer = reflect.TypeOf(&serializer).Elem()
}

// SerializerSwitch defines switch of BinarySerializer check
type SerializerSwitch byte

const (
	SerializerDisable    SerializerSwitch = iota // disable Serializer
	SerializerCheck                              // enable Serializer but need check
	SerializerCheckFalse                         // enable and do not need check,result false
	SerializerCheckOk                            // enable and do not need check,result true
)

// String return name of this switch
func (ss SerializerSwitch) String() string {
	switch ss {
	case SerializerDisable:
		return "SerializerDisable"
	case SerializerCheck:
		return "SerializerCheck"
	case SerializerCheckOk:
		return "SerializerCheckOk"
	}
	panic(fmt.Errorf("SerializerUnknown"))
}

// Enable returns if BinarySerializer check is enable
func (ss SerializerSwitch) Enable() bool {
	return ss >= SerializerCheck
}

// NeedCheck returns if need check BinarySerializer
func (ss SerializerSwitch) NeedCheck() bool {
	return ss == SerializerCheck
}

// CheckFail returns if can use BinarySerializer directly
func (ss SerializerSwitch) CheckFalse() bool {
	return ss == SerializerCheckFalse
}

// NeedCheck returns if can use BinarySerializer directly
func (ss SerializerSwitch) CheckOk() bool {
	return ss == SerializerCheckOk
}

//// Check returns if t is a BinarySerializer when enable
//func (ss SerializerSwitch) Check(t reflect.Type) bool {
//	switch {
//	case ss.CheckOk() || ss.NeedCheck() && querySerializer(indirectType(t)):
//		return true
//	case !ss.Enable() || ss.CheckFalse():
//		fallthrough
//	default:
//		return false
//	}
//}

// SubSwitch returns SerializerSwitch for sub-data of struct/map/slice/array
func (ss SerializerSwitch) SubSwitchCheck(t reflect.Type) SerializerSwitch {
	if !ss.Enable() {
		return SerializerDisable
	}
	return ss.subSwitch(querySerializer(indirectType(t)))
}

func (ss SerializerSwitch) subSwitch(isSerializer bool) SerializerSwitch {
	if !ss.Enable() {
		return SerializerDisable
	}
	if isSerializer {
		return SerializerCheckOk
	}
	return SerializerCheckFalse
}

func toplvSerializer(enable bool) SerializerSwitch {
	if enable {
		return SerializerCheck
	}
	return SerializerDisable
}

//CheckSerializer check if t implements BinarySerializer
func CheckSerializer(x interface{}) bool {
	return querySerializer(indirectType(reflect.TypeOf(x)))
}

//CheckSerializerDeep check if t or *t implements BinarySerializer
func CheckSerializerDeep(t reflect.Type) bool {
	return querySerializer(indirectType(t))
}
