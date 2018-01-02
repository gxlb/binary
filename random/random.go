// package random provides rander to generate rand data.
package random

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

const (
	defaultStringLen = 100
)

var (
	strFull     = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890_-#@!$&")
	defaultRand = NewRand(0)
)

// Default return default Rand object
func Default() *Rand {
	return defaultRand
}

// RandSeed generate a seed for rand
func RandSeed(init uint64) uint64 {
	if 0 == init {
		init = uint64(time.Now().Unix())
	}
	rnd := init*0xA7A263F04949875D + 0x88AB71F41758C2AF
	return rnd
}

// NewRand create a new random object with init
func NewRand(init uint64) *Rand {
	return &Rand{seed: RandSeed(init)}
}

// Rand generate random data.
type Rand struct {
	seed uint64
}

// Rand generate next rand number
func (rnd *Rand) Rand() uint64 {
	n := rnd.seed*0x9A7A8C6232996B69 + 0xFD06D0380857153B
	rnd.seed = n
	return n
}

// RandRange generate a 32-bit rand number with range.
func (rnd *Rand) RandRange(min, max uint32) uint32 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Uint32()
	ret := r%d + min

	return ret
}

// RandRange64 generate a 64-bit rand number with range.
func (rnd *Rand) RandRange64(min, max uint64) uint64 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Uint64()
	ret := r%d + min

	return ret
}

// RandMax generate a 32-bit rand number with max value.
func (rnd *Rand) RandMax(max uint32) uint32 {
	return rnd.RandRange(0, max-1)
}

// RandMax generate a 64-bit rand number with max value.
func (rnd *Rand) RandMax64(max uint64) uint64 {
	return rnd.RandRange64(0, max-1)
}

// Seed returns current seed.
func (rnd *Rand) Seed() uint64 {
	return rnd.seed
}

// CopyNew create a new Rand pointer from rnd.
// It will call rnd.Rand before copy.
func (rnd *Rand) CopyNew() *Rand {
	rnd.Rand() //next rand number
	return &Rand{seed: rnd.seed}
}

// Copy create a copy of Rand object from rnd.
// It will call rnd.Rand before copy.
func (rnd *Rand) Copy() Rand {
	rnd.Rand() //next rand number
	return *rnd
}

// Reset set seed with a random seed.
func (rnd *Rand) Reset() uint64 {
	return rnd.Srand(RandSeed(0))
}

// Srand set seed of the object.
func (rnd *Rand) Srand(seed uint64) uint64 {
	ret := rnd.seed
	rnd.seed = seed
	return ret
}

// String generate a random string value with length.
func (rnd *Rand) String(length int) string {
	b := make([]byte, length, length)
	for i := 0; i < length; i++ {
		b[i] = strFull[rnd.RandMax(uint32(len(strFull)))]
	}
	return string(b)
}

// Bool generate a random bool value.
func (rnd *Rand) Bool() bool {
	return rnd.Uint8()&0x1 == 0
}

// Uint generate a random uint value.
func (rnd *Rand) Uint() uint {
	return uint(rnd.Uint64())
}

// Int generate a random int value.
func (rnd *Rand) Int() int {
	return int(rnd.Uint())
}

// Uint8 generate a random uint8 value.
func (rnd *Rand) Uint8() uint8 {
	return uint8(rnd.Rand() >> 54 & 0xFF)
}

// Int8 generate a random int8 value.
func (rnd *Rand) Int8() int8 {
	return int8(rnd.Uint8())
}

// Uint16 generate a random uint16 value.
func (rnd *Rand) Uint16() uint16 {
	return uint16(rnd.Rand() >> 46 & 0xFFFF)
}

// Int16 generate a random int16 value.
func (rnd *Rand) Int16() int16 {
	return int16(rnd.Uint16())
}

// Uint32 generate a random uint32 value.
func (rnd *Rand) Uint32() uint32 {
	return uint32(rnd.Rand() >> 30 & 0xFFFFFFFF)
}

// Int32 generate a random int32 value.
func (rnd *Rand) Int32() int32 {
	return int32(rnd.Uint32())
}

// Uint64 generate a random uint64 value.
func (rnd *Rand) Uint64() uint64 {
	v := uint64(rnd.Uint32())<<32 + uint64(rnd.Uint32())
	return v
}

// Int64 generate a random int64 value.
func (rnd *Rand) Int64() int64 {
	return int64(rnd.Uint64())
}

// Float32 generate a random float32 value.
func (rnd *Rand) Float32() float32 {
	return math.Float32frombits(rnd.Uint32())
}

// Float64 generate a random float64 value.
func (rnd *Rand) Float64() float64 {
	return math.Float64frombits(rnd.Uint64())
}

// Complex64 generate a random complex64 value.
func (rnd *Rand) Complex64() complex64 {
	r := rnd.Float32()
	i := rnd.Float32()
	return complex(r, i)
}

// Complex128 generate a random complex128 value.
func (rnd *Rand) Complex128() complex128 {
	r := rnd.Float64()
	i := rnd.Float64()
	return complex(r, i)
}

// Value writes x with random data.
// x must be a pointer value.
func (rnd *Rand) Value(x interface{}) error {
	return rnd.ValueX(x, 0, 0, 0)
}

// ValueX writes x with random data.
// x must be a pointer value.
// if seed > 0, it will generate by a copy of rnd with seed.
// minLen, maxLen represents length range of slice/map.
func (rnd *Rand) ValueX(x interface{}, seed uint64, minLen, maxLen uint32) error {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("can only set rand value by non-nil pointer, got %s", v.Type().String())
	}
	r := rnd
	if seed != 0 {
		c := rnd.Copy()
		c.Srand(seed)
		r = &c
	}
	return r.value(v.Elem(), minLen, maxLen)
}

func (rnd *Rand) value(v reflect.Value, minLen, maxLen uint32) error {
	switch k := v.Kind(); k {
	case reflect.Int:
		v.SetInt(int64(rnd.Int()))
	case reflect.Uint:
		v.SetUint(uint64(rnd.Uint()))
	case reflect.Bool:
		v.SetBool(rnd.Bool())
	case reflect.Int8:
		v.SetInt(int64(rnd.Int8()))
	case reflect.Int16:
		v.SetInt(int64(rnd.Int16()))
	case reflect.Int32:
		v.SetInt(int64(rnd.Int32()))
	case reflect.Int64:
		v.SetInt(int64(rnd.Int64()))
	case reflect.Uint8:
		v.SetUint(uint64(rnd.Uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(rnd.Uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(rnd.Uint32()))
	case reflect.Uint64:
		v.SetUint(uint64(rnd.Uint64()))
	case reflect.Float32:
		v.SetFloat(float64(rnd.Float32()))
	case reflect.Float64:
		v.SetFloat(rnd.Float64())
	case reflect.Complex64:
		v.SetComplex(complex128(rnd.Complex64()))
	case reflect.Complex128:
		v.SetComplex(rnd.Complex128())
	case reflect.String:
		v.SetString(rnd.String(rnd.length(defaultStringLen, 0)))

	case reflect.Slice, reflect.Array:
		if k == reflect.Slice {
			length := rnd.length(minLen, maxLen)
			v.Set(reflect.MakeSlice(v.Type(), length, length))
		}
		for i := 0; i < v.Len(); i++ {
			rnd.value(v.Index(i).Addr(), minLen, maxLen)
		}

	case reflect.Map:
		length := rnd.length(minLen, maxLen)
		v.Set(reflect.MakeMap(v.Type()))
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()
		for i := 0; i < length; i++ {
			key := reflect.New(kt).Elem()
			val := reflect.New(vt).Elem()
			rnd.value(key, minLen, maxLen)
			rnd.value(val, minLen, maxLen)
			v.SetMapIndex(key, val)
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			rnd.value(f, minLen, maxLen)
		}

	case reflect.Ptr:
		if v.IsNil() {
			elemT := v.Type().Elem()
			v.Set(reflect.New(elemT))
		}
		return rnd.value(v.Elem(), minLen, maxLen)

	default:
		return fmt.Errorf("random.Value: unsupported type [%s]", v.Type().String())
	}

	return nil
}

// length generate random lenth with range [minLen, maxLen]
func (rnd *Rand) length(minLen, maxLen uint32) int {
	if minLen == 0 {
		minLen = 10
	}
	if maxLen == 0 {
		maxLen = minLen + minLen/2
	}
	return int(rnd.RandRange(minLen, maxLen))
}

// New create a new object with type of x.
func New(x interface{}, length int) interface{} {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_new := reflect.New(t)
	switch t.Kind() {
	case reflect.Slice:
		_new.Set(reflect.MakeSlice(t, t.Len(), t.Len()))
	case reflect.Map:
		_new.Set(reflect.MakeMap(t))
	}
	return _new.Interface()
}

// deepNew create a new object with type of v.
func deepNew(v reflect.Value, length int) {
	t := v.Type()
	_new := reflect.New(t)
	k := t.Kind()
	switch k {
	case reflect.Slice:
		_new.Set(reflect.MakeSlice(t, length, length))
	case reflect.Map:
		_new.Set(reflect.MakeMap(t))
	}
	v.Set(_new)
	if k == reflect.Ptr {
		deepNew(v.Elem(), length)
	}
}
