//rand data generator

package random

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

var (
	strFull = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-#@!$&")
	rand    = NewRand(1)
)

//generate a seed for rand
func RandSeed(init uint32) uint32 {
	if 0 == init {
		init = uint32(time.Now().Unix())
	}
	rnd := init*4294967291 + 8615693
	return rnd
}

func NewRand(init uint32) *Rand {
	return &Rand{seed: RandSeed(init)}
}

type Rand struct {
	seed uint32
}

//next rand number
func (rnd *Rand) Rand() uint32 {
	n := rnd.seed*7368787 + 2750159
	rnd.seed = n
	return n
}

//generate rand number in range
func (rnd *Rand) RandRange(min, max uint32) uint32 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Rand()
	ret := r%d + min

	return ret
}

//generate rand number in range
func (rnd *Rand) RandRange64(min, max uint64) uint64 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Uint64()
	ret := r%d + min

	return ret
}

//generate rand number with max value
func (rnd *Rand) RandMax(max uint32) uint32 {
	return rnd.RandRange(0, max-1)
}

//generate rand number with max value
func (rnd *Rand) RandMax64(max uint64) uint64 {
	return rnd.RandRange64(0, max-1)
}

//get seed
func (rnd *Rand) Seed() uint32 {
	return rnd.seed
}

func (rnd *Rand) CopyNew() *Rand {
	return &Rand{seed: rnd.seed}
}

func (rnd *Rand) Copy() Rand {
	return *rnd
}

//set seed
func (rnd *Rand) Srand(seed uint32) uint32 {
	ret := rnd.seed
	rnd.seed = seed
	return ret
}

func (rnd *Rand) String(length int) string {
	b := make([]byte, length, length)
	for i := 0; i < length; i++ {
		b[i] = strFull[rnd.RandMax(uint32(len(strFull)))]
	}
	return string(b)
}

func (rnd *Rand) Bool() bool {
	return rnd.Rand()&0x1 == 0
}

func (rnd *Rand) Uint() uint {
	return uint(rnd.Uint64())
}
func (rnd *Rand) Int() int {
	return int(rnd.Uint())
}
func (rnd *Rand) Uint8() uint8 {
	return uint8(rnd.Rand() & 0xFF)
}
func (rnd *Rand) Int8() int8 {
	return int8(rnd.Uint8())
}

func (rnd *Rand) Uint16() uint16 {
	v := uint16(0)
	for i := 0; i < 2; i++ {
		v = v<<8 + uint16(rnd.Uint8())
	}
	return v
}

func (rnd *Rand) Int16() int16 {
	return int16(rnd.Uint16())
}

func (rnd *Rand) Uint32() uint32 {
	v := uint32(0)
	for i := 0; i < 4; i++ {
		v = v<<8 + uint32(rnd.Uint8())
	}
	return v
}

func (rnd *Rand) Int32() int32 {
	return int32(rnd.Uint32())
}

func (rnd *Rand) Uint64() uint64 {
	v := uint64(0)
	for i := 0; i < 8; i++ {
		v = v<<8 + uint64(rnd.Uint8())
	}
	return v
}

func (rnd *Rand) Int64() int64 {
	return int64(rnd.Uint64())
}

func (rnd *Rand) Float32() float32 {
	return math.Float32frombits(rnd.Uint32())
}

func (rnd *Rand) Float64() float64 {
	return math.Float64frombits(rnd.Uint64())
}

func (rnd *Rand) Complex64() complex64 {
	r := rnd.Float32()
	i := rnd.Float32()
	return complex(r, i)
}

func (rnd *Rand) Complex128() complex128 {
	r := rnd.Float64()
	i := rnd.Float64()
	return complex(r, i)
}

//generate rand value for x
func (rnd *Rand) Value(x interface{}) error {
	return rnd.ValueX(x, 0, 0, 0, 0, 0)
}

func (rnd *Rand) ValueX(x interface{}, seed uint32, minLen, maxLen uint32, min, max int64) error {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("can only set rand value by non-nil pointer, got %s", v.Type().String())
	}
	return rnd.value(v.Elem(), minLen, maxLen, min, max)
}

func (rnd *Rand) value(v reflect.Value, minLen, maxLen uint32, min, max int64) error {
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
