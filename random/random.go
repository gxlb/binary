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
	rand    = NewRand(0)
)

//generate a seed for rand
func RandSeed(init uint64) uint64 {
	if 0 == init {
		init = uint64(time.Now().Unix())
	}
	rnd := init*0xA7A263F04949875D + 0x88AB71F41758C2AF
	return rnd
}

func NewRand(init uint64) *Rand {
	return &Rand{seed: RandSeed(init)}
}

type Rand struct {
	seed uint64
}

//next rand number
func (rnd *Rand) Rand() uint64 {
	n := rnd.seed*0x9A7A8C6232996B69 + 0xFD06D0380857153B
	rnd.seed = n
	return n
}

//generate rand number in range
func (rnd *Rand) RandRange(min, max uint32) uint32 {
	if max < min {
		max, min = min, max
	}
	d := max - min + 1
	r := rnd.Uint32()
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
func (rnd *Rand) Seed() uint64 {
	return rnd.seed
}

func (rnd *Rand) CopyNew() *Rand {
	return &Rand{seed: rnd.seed}
}

func (rnd *Rand) Copy() Rand {
	return *rnd
}

//reset generate a random seed
func (rnd *Rand) Reset() uint64 {
	return rnd.Srand(RandSeed(0))
}

//set seed
func (rnd *Rand) Srand(seed uint64) uint64 {
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
	return rnd.Uint8()&0x1 == 0
}

func (rnd *Rand) Uint() uint {
	return uint(rnd.Uint64())
}
func (rnd *Rand) Int() int {
	return int(rnd.Uint())
}
func (rnd *Rand) Uint8() uint8 {
	return uint8(rnd.Rand() >> 54 & 0xFF)
}
func (rnd *Rand) Int8() int8 {
	return int8(rnd.Uint8())
}

func (rnd *Rand) Uint16() uint16 {
	return uint16(rnd.Rand() >> 46 & 0xFFFF)
}

func (rnd *Rand) Int16() int16 {
	return int16(rnd.Uint16())
}

func (rnd *Rand) Uint32() uint32 {
	return uint32(rnd.Rand() >> 30 & 0xFFFFFFFF)
}

func (rnd *Rand) Int32() int32 {
	return int32(rnd.Uint32())
}

func (rnd *Rand) Uint64() uint64 {
	v := uint64(rnd.Uint32())<<32 + uint64(rnd.Uint32())
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
	return rnd.ValueX(x, 0, 0, 0)
}

func (rnd *Rand) ValueX(x interface{}, seed uint64, minLen, maxLen uint32) error {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("can only set rand value by non-nil pointer, got %s", v.Type().String())
	}
	if seed != 0 {
		rnd.Srand(seed)
	}
	if 0 == minLen {
		minLen = 10
	}
	return rnd.value(v.Elem(), minLen, maxLen)
}

func (rnd *Rand) value(v reflect.Value, minLen, maxLen uint32) error {
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
		v.Set(reflect.ValueOf(rnd.String(rnd.length(minLen, maxLen))))

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

func (rnd *Rand) length(minLen, maxLen uint32) int {
	if minLen == 0 {
		minLen = 10
	}
	if maxLen == 0 {
		maxLen = minLen + minLen/2
	}
	return int(rnd.RandRange(minLen, maxLen))
}

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
