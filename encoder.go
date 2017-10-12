package binary

import (
	"fmt"
	"math"
	"reflect"
)

func NewEncoder(size int) *Encoder {
	return NewEncoderEndian(size, DefaultEndian)
}

func NewEncoderEndian(size int, endian Endian) *Encoder {
	p := &Encoder{}
	p.Init(size, endian)
	return p
}

type Encoder struct {
	coder
}

func (this *Encoder) Init(size int, endian Endian) {
	this.buff = make([]byte, size)
	this.pos = 0
	this.endian = endian
}

func (this *Encoder) Bool(x bool) {
	b := this.reserve(1)
	if x {
		b[0] = 1
	} else {
		b[0] = 0
	}
}

func (this *Encoder) Int8(x int8) {
	this.Uint8(uint8(x))
}

func (this *Encoder) Uint8(x uint8) {
	b := this.reserve(1)
	b[0] = x
}

func (this *Encoder) Int16(x int16) {
	this.Uint16(uint16(x))
}

func (this *Encoder) Uint16(x uint16) {
	b := this.reserve(2)
	this.endian.PutUint16(b, x)
}

func (this *Encoder) Int32(x int32) {
	this.Uint32(uint32(x))
}

func (this *Encoder) Uint32(x uint32) {
	b := this.reserve(4)
	this.endian.PutUint32(b, x)
}

func (this *Encoder) Int64(x int64) {
	this.Uint64(uint64(x))
}

func (this *Encoder) Uint64(x uint64) {
	b := this.reserve(8)
	this.endian.PutUint64(b, x)
}

func (this *Encoder) Float32(x float32) {
	this.Uint32(math.Float32bits(x))
}

func (this *Encoder) Float64(x float64) {
	this.Uint64(math.Float64bits(x))
}

func (this *Encoder) Complex64(x complex64) {
	this.Uint32(math.Float32bits(real(x)))
	this.Uint32(math.Float32bits(imag(x)))
}

func (this *Encoder) Complex128(x complex128) {
	this.Uint64(math.Float64bits(real(x)))
	this.Uint64(math.Float64bits(imag(x)))
}

func (this *Encoder) String(x string) {
	_b := []byte(x)
	size := len(_b)
	this.Uint32(uint32(size))
	buff := this.reserve(size)
	copy(buff, _b)
}

func (this *Encoder) Value(x interface{}) error {
	if this.fastValue(x) { //fast value path
		return nil
	}
	v := reflect.ValueOf(x)
	return this.value(reflect.Indirect(v))
}

func (this *Encoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case bool:
		this.Bool(d)
	case int8:
		this.Int8(d)
	case uint8:
		this.Uint8(d)
	case int16:
		this.Int16(d)
	case uint16:
		this.Uint16(d)
	case int32:
		this.Int32(d)
	case uint32:
		this.Uint32(d)
	case float32:
		this.Float32(d)
	case int64:
		this.Int64(d)
	case uint64:
		this.Uint64(d)
	case float64:
		this.Float64(d)
	case complex64:
		this.Complex64(d)
	case complex128:
		this.Complex128(d)
	case string:
		this.String(d)
	case []bool:
		l := len(d)
		this.Uint32(__cntType(l))
		var b []byte
		for i := 0; i < l; i++ {
			bit := i % 8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b = this.reserve(1)
				b[0] = 0
			}
			if x := d[i]; x {
				b[0] |= mask
			}
		}

	case []int8:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Int8(d[i])
		}
	case []uint8:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Uint8(d[i])
		}
	case []int16:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Int16(d[i])
		}
	case []uint16:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Uint16(d[i])
		}

	case []int32:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Int32(d[i])
		}
	case []uint32:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Uint32(d[i])
		}
	case []int64:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Int64(d[i])
		}
	case []uint64:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Uint64(d[i])
		}
	case []float32:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Float32(d[i])
		}
	case []float64:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Float64(d[i])
		}
	case []complex64:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Complex64(d[i])
		}
	case []complex128:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.Complex128(d[i])
		}
	case []string:
		l := len(d)
		this.Uint32(__cntType(len(d)))
		for i := 0; i < l; i++ {
			this.String(d[i])
		}
		//	case *bool:
		//		this.Bool(*d)
		//	case *int8:
		//		this.Int8(*d)
		//	case *uint8:
		//		this.Uint8(*d)
		//	case *int16:
		//		this.Int16(*d)
		//	case *uint16:
		//		this.Uint16(*d)
		//	case *int32:
		//		this.Int32(*d)
		//	case *uint32:
		//		this.Uint32(*d)
		//	case *float32:
		//		this.Float32(*d)
		//	case *int64:
		//		this.Int64(*d)
		//	case *uint64:
		//		this.Uint64(*d)
		//	case *float64:
		//		this.Float64(*d)
		//	case *complex64:
		//		this.Complex64(*d)
		//	case *complex128:
		//		this.Complex128(*d)
		//	case *string:
		//		this.String(*d)
	default:
		return false
	}
	return true

}

func (this *Encoder) value(v reflect.Value) error {
	//	defer func() {
	//		fmt.Printf("Encoder:after value(%#v)=%d\n", v.Interface(), this.pos)
	//	}()
	switch k := v.Kind(); k {
	case reflect.Bool:
		this.Bool(v.Bool())

	case reflect.Int8:
		this.Int8(int8(v.Int()))
	case reflect.Int16:
		this.Int16(int16(v.Int()))
	case reflect.Int32:
		this.Int32(int32(v.Int()))
	case reflect.Int, reflect.Int64:
		this.Int64(v.Int())

	case reflect.Uint8:
		this.Uint8(uint8(v.Uint()))
	case reflect.Uint16:
		this.Uint16(uint16(v.Uint()))
	case reflect.Uint32:
		this.Uint32(uint32(v.Uint()))
	case reflect.Uint, reflect.Uint64:
		this.Uint64(v.Uint())

	case reflect.Float32:
		this.Uint32(math.Float32bits(float32(v.Float())))
	case reflect.Float64:
		this.Uint64(math.Float64bits(v.Float()))

	case reflect.Complex64:
		x := v.Complex()
		this.Uint32(math.Float32bits(float32(real(x))))
		this.Uint32(math.Float32bits(float32(imag(x))))
	case reflect.Complex128:
		x := v.Complex()
		this.Uint64(math.Float64bits(real(x)))
		this.Uint64(math.Float64bits(imag(x)))

	case reflect.String:
		this.String(v.String())

	case reflect.Slice, reflect.Array:
		if this.boolArray(v) < 0 { //deal with bool array first
			l := v.Len()
			this.Uint32(__cntType(l))
			for i := 0; i < l; i++ {
				this.value(v.Index(i))
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		l := len(keys)
		this.Uint32(__cntType(l))
		for i := 0; i < l; i++ {
			key := keys[i]
			this.value(key)
			this.value(v.MapIndex(key))
		}
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// see comment for corresponding code in decoder.value()
			if f := v.Field(i); validField(f, t.Field(i)) {
				this.value(f)
			} else {
				//this.Skip(sizeofEmptyValue(f))
			}
		}
	case reflect.Ptr:
		if !v.IsNil() {
			if e := v.Elem(); e.Kind() != reflect.Ptr {
				return this.value(e)
			}
		} else {
			this.Skip(sizeofEmptyValue(v))
		}
	default:
		return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Type().String())
	}
	return nil
}

// encode bool array
func (this *Encoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			l := v.Len()
			this.Uint32(__cntType(l))
			var b []byte
			for i := 0; i < l; i++ {
				bit := i % 8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = this.reserve(1)
					b[0] = 0
					//fmt.Println("Encoder.boolArray", i, bit, this.pos)
				}
				if x := v.Index(i).Bool(); x {
					b[0] |= mask
				}
			}
			return sizeofBoolArray(l)
		}
	}
	return -1
}
