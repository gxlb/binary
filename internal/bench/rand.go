//rand data generator

package bench

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

const (
	rndA = 7368787
	rndC = 2750159
)

var (
	strFull = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-#@!$&")
	rand    = NewRand64(1)
)

//generate a seed for rand
func RandSeed32(init uint32) uint32 {
	if 0 == init {
		init = uint32(time.Now().Unix())
	}
	rnd := init*4294967291 + 8615693
	rnd = rnd*4294967291 + 8615693
	return rnd
}

type Rand32 struct {
	seed uint32
}

func NewRand64(init uint32) *Rand32 {
	return &Rand32{seed: RandSeed32(init)}
}

//next rand number
func (rnd *Rand32) Rand() uint32 {
	n := rnd.seed*rndA + rndC
	rnd.seed = n
	return n
}

//generate rand number in range
func (rnd *Rand32) RandRange(min, max uint32) uint32 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Rand()
	ret := r%d + min

	return ret
}

//generate rand number with max value
func (rnd *Rand32) RandMax(max uint32) uint32 {
	return rnd.RandRange(0, max-1)
}

//get seed
func (rnd *Rand32) Seed() uint32 {
	return rnd.seed
}

//set seed
func (rnd *Rand32) Srand(seed uint32) uint32 {
	ret := rnd.seed
	rnd.seed = seed
	return ret
}

func (rnd *Rand32) String(length int) string {
	b := make([]byte, length, length)
	for i := 0; i < length; i++ {
		b[i] = strFull[rnd.RandMax(uint32(len(strFull)))]
	}
	return string(b)
}

func (rnd *Rand32) Bool() bool {
	return rnd.Rand()&0x1 == 0
}

func (rnd *Rand32) Uint() uint {
	return uint(rnd.Uint64())
}
func (rnd *Rand32) Int() int {
	return int(rnd.Uint())
}
func (rnd *Rand32) Uint8() uint8 {
	return uint8(rnd.Rand() & 0xFF)
}
func (rnd *Rand32) Int8() int8 {
	return int8(rnd.Uint8())
}
func (rnd *Rand32) Uint16() uint16 {
	v := uint16(0)
	for i := 0; i < 2; i++ {
		v = v<<8 + uint16(rnd.Uint8())
	}
	return v
}
func (rnd *Rand32) Int16() int16 {
	return int16(rnd.Uint16())
}
func (rnd *Rand32) Uint32() uint32 {
	v := uint32(0)
	for i := 0; i < 4; i++ {
		v = v<<8 + uint32(rnd.Uint8())
	}
	return v
}
func (rnd *Rand32) Int32() int32 {
	return int32(rnd.Uint32())
}
func (rnd *Rand32) Uint64() uint64 {
	v := uint64(0)
	for i := 0; i < 8; i++ {
		v = v<<8 + uint64(rnd.Uint8())
	}
	return v
}
func (rnd *Rand32) Int64() int64 {
	return int64(rnd.Uint64())
}
func (rnd *Rand32) Float32() float32 {
	return math.Float32frombits(rnd.Uint32())
}
func (rnd *Rand32) Float64() float64 {
	return math.Float64frombits(rnd.Uint64())
}
func (rnd *Rand32) Complex64() complex64 {
	r := rnd.Float32()
	i := rnd.Float32()
	return complex(r, i)
}
func (rnd *Rand32) Complex128() complex128 {
	r := rnd.Float64()
	i := rnd.Float64()
	return complex(r, i)
}

//generate rand value for x
func (rnd *Rand32) Value(x interface{}) error {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("can only set rand value by non-nil pointer, got %s", v.Type().String())
	}
	return rnd.value(v.Elem())
}

func (rnd *Rand32) value(v reflect.Value) error {
	switch k := v.Kind(); k {
	case reflect.Int:
		v.Set(reflect.ValueOf(rnd.Int()))
	case reflect.Uint:
		v.Set(reflect.ValueOf(rnd.Uint()))
	case reflect.Bool:
		v.Set(reflect.ValueOf(rnd.Bool()))
	case reflect.Int8:
		v.Set(reflect.ValueOf(rnd.Int8()))
	case reflect.Int16:
		v.Set(reflect.ValueOf(rnd.Int16()))
	case reflect.Int32:
		v.Set(reflect.ValueOf(rnd.Int32()))
	case reflect.Int64:
		v.Set(reflect.ValueOf(rnd.Int64()))
	case reflect.Uint8:
		v.Set(reflect.ValueOf(rnd.Uint8()))
	case reflect.Uint16:
		v.Set(reflect.ValueOf(rnd.Uint16()))
	case reflect.Uint32:
		v.Set(reflect.ValueOf(rnd.Uint32()))
	case reflect.Uint64:
		v.Set(reflect.ValueOf(rnd.Uint64()))
	case reflect.Float32:
		v.Set(reflect.ValueOf(rnd.Float32()))
	case reflect.Float64:
		v.Set(reflect.ValueOf(rnd.Float64()))
	case reflect.Complex64:
		v.Set(reflect.ValueOf(rnd.Complex64()))
	case reflect.Complex128:
		v.Set(reflect.ValueOf(rnd.Complex128()))
	case reflect.String:
		v.Set(reflect.ValueOf(rnd.String(int(rnd.RandMax(100)))))

	case reflect.Slice, reflect.Array:

	case reflect.Map:

	case reflect.Struct:

	case reflect.Ptr:

	default:
		//return typeError("binary.Encoder.Value: unsupported type [%s]", v.Type(), true)

	}
	return nil
}
