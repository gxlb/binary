package binary

import (
	"fmt"
	"math"
	"reflect"
)

func NewDecoder(buff []byte) *Decoder {
	return NewDecoderEndian(buff, DefaultEndian)
}

func NewDecoderEndian(buff []byte, endian Endian) *Decoder {
	p := &Decoder{}
	p.Init(buff, endian)
	return p
}

type Decoder struct {
	coder
}

func (this *Decoder) Init(buff []byte, endian Endian) {
	this.buff = buff
	this.pos = 0
	this.endian = endian
}

func (this *Decoder) Bool() bool {
	b := this.reserve(1)
	x := b[0]
	return x != 0
}

func (this *Decoder) Int8() int8 {
	return int8(this.Uint8())
}

func (this *Decoder) Uint8() uint8 {
	b := this.reserve(1)
	x := b[0]
	return x
}

func (this *Decoder) Int16() int16 {
	return int16(this.Uint16())
}

func (this *Decoder) Uint16() uint16 {
	b := this.reserve(2)
	x := this.endian.Uint16(b)
	return x
}

func (this *Decoder) Int32() int32 {
	return int32(this.Uint32())
}

func (this *Decoder) Uint32() uint32 {
	b := this.reserve(4)
	x := this.endian.Uint32(b)
	return x
}

func (this *Decoder) Int64() int64 {
	return int64(this.Uint64())
}

func (this *Decoder) Uint64() uint64 {
	b := this.reserve(8)
	x := this.endian.Uint64(b)
	return x
}

func (this *Decoder) Float32() float32 {
	x := math.Float32frombits(this.Uint32())
	return x
}

func (this *Decoder) Float64() float64 {
	x := math.Float64frombits(this.Uint64())
	return x
}

func (this *Decoder) Complex64() complex64 {
	x := complex(this.Float32(), this.Float32())
	return x
}

func (this *Decoder) Complex128() complex128 {
	x := complex(this.Float64(), this.Float64())
	return x
}

func (this *Decoder) String() string {
	var s __cntType = this.Uint32()
	size := int(s)
	b := this.reserve(size)
	return string(b)
}

func (this *Decoder) Value(x interface{}) error {
	if this.fastValue(x) { //fast value path
		return nil
	}

	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr { //only support decode for pointer interface
		return this.value(reflect.Indirect(v))
	} else {
		return fmt.Errorf("binary.Decoder.Value: non-pointer type [%s]", v.Type().String())
	}
}

//func (this *Decoder) getIntValue(kind reflect.Kind) uint64 {
//	v := uint64(0)
//	switch kind {
//	case reflect.Int8:
//		v = uint64(this.Int8())
//	case reflect.Int16:
//		v = uint64(this.Int16())
//	case reflect.Int32:
//		v = uint64(this.Int32())
//	case reflect.Int64:
//		v = uint64(this.Int64())

//	case reflect.Uint8:
//		v = uint64(this.Uint8())
//	case reflect.Uint16:
//		v = uint64(this.Uint16())
//	case reflect.Uint32:
//		v = uint64(this.Uint32())
//	case reflect.Uint64:
//		v = this.Uint64()
//	default:
//		panic(kind)
//	}
//	return v
//}

func (this *Decoder) value(v reflect.Value) error {
	//	defer func() {
	//		fmt.Printf("Decoder:after value(%#v)=%d\n", v.Interface(), this.pos)
	//	}()
	switch k := v.Kind(); k {
	case reflect.Bool:
		v.SetBool(this.Bool())

	case reflect.Int8:
		v.SetInt(int64(this.Int8()))
	case reflect.Int16:
		v.SetInt(int64(this.Int16()))
	case reflect.Int32:
		v.SetInt(int64(this.Int32()))
	case reflect.Int, reflect.Int64:
		v.SetInt(this.Int64())

	case reflect.Uint8:
		v.SetUint(uint64(this.Uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(this.Uint16()))
	case reflect.Uint32:
		v.SetUint(uint64(this.Uint32()))
	case reflect.Uint, reflect.Uint64:
		v.SetUint(this.Uint64())

	case reflect.Float32:
		v.SetFloat(float64(this.Float32()))
	case reflect.Float64:
		v.SetFloat(this.Float64())

	case reflect.Complex64:
		v.SetComplex(complex(
			float64(this.Float32()),
			float64(this.Float32()),
		))
	case reflect.Complex128:
		v.SetComplex(complex(
			this.Float64(),
			this.Float64(),
		))
	case reflect.String:
		v.SetString(this.String())

	case reflect.Slice, reflect.Array:
		if this.boolArray(v) < 0 { //deal with bool array first
			var s __cntType = this.Uint32()
			size := int(s)
			if k == reflect.Slice { //make a new slice
				ns := reflect.MakeSlice(v.Type(), size, size)
				v.Set(ns)
			}

			l := v.Len()
			for i := 0; i < size; i++ {
				if i < l {
					this.value(v.Index(i))
				} else {
					this.Skip(this.sizeofType(v.Type().Elem()))
				}
			}
		}
	case reflect.Map:
		var s __cntType = this.Uint32()
		size := int(s)
		newmap := reflect.MakeMap(v.Type())
		v.Set(newmap)
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()

		for i := 0; i < size; i++ {
			//fmt.Printf("key:%#v value:%#v\n", key.Elem().Interface(), value.Elem().Interface())
			key := reflect.New(kt).Elem()
			value := reflect.New(vt).Elem()
			this.value(key)
			this.value(value)
			v.SetMapIndex(key, value)
		}
	case reflect.Struct:
		t := v.Type()
		l := v.NumField()
		for i := 0; i < l; i++ {
			// Note: Calling v.CanSet() below is an optimization.
			// It would be sufficient to check the field name,
			// but creating the StructField info for each field is
			// costly (run "go test -bench=ReadStruct" and compare
			// results when making changes to this code).
			if f := v.Field(i); validField(f, t.Field(i)) {
				//fmt.Printf("field(%d) [%s] \n", i, t.Field(i).Name)
				this.value(f)
			} else {
				//this.Skip(this.sizeofType(f.Type()))
			}
		}
	default:
		if newPtr(v) {
			return this.value(v.Elem())
		} else {
			return fmt.Errorf("binary.Decoder.Value: unsupported type [%s]", v.Type().String())
		}
	}
	return nil
}

func (this *Decoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case *bool:
		*d = this.Bool()
	case *int8:
		*d = this.Int8()
	case *uint8:
		*d = this.Uint8()

	case *int16:
		*d = this.Int16()
	case *uint16:
		*d = this.Uint16()

	case *int32:
		*d = this.Int32()
	case *uint32:
		*d = this.Uint32()
	case *float32:
		*d = this.Float32()

	case *int64:
		*d = this.Int64()
	case *uint64:
		*d = this.Uint64()
	case *float64:
		*d = this.Float64()
	case *complex64:
		*d = this.Complex64()

	case *complex128:
		*d = this.Complex128()

	case *string:
		*d = this.String()

	case *[]bool:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]bool, l)
		var b []byte
		for i := 0; i < l; i++ {
			_, bit := i/8, i%8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b = this.reserve(1)
			}
			x := ((b[0] & mask) != 0)
			(*d)[i] = x
		}

	case *[]int8:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]int8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int8()
		}
	case *[]uint8:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]uint8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint8()
		}
	case *[]int16:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]int16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int16()
		}
	case *[]uint16:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]uint16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint16()
		}
	case *[]int32:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]int32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int32()
		}
	case *[]uint32:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]uint32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint32()
		}
	case *[]int64:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]int64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Int64()
		}
	case *[]uint64:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]uint64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Uint64()
		}
	case *[]float32:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]float32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Float32()
		}
	case *[]float64:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]float64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Float64()
		}
	case *[]complex64:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]complex64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Complex64()
		}
	case *[]complex128:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]complex128, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.Complex128()
		}
	case *[]string:
		var s __cntType = this.Uint32()
		l := int(s)
		*d = make([]string, l)
		for i := 0; i < l; i++ {
			(*d)[i] = this.String()
		}
	default:
		return false
	}
	return true
}

//get size of specific type from current buffer
func (this *Decoder) sizeofType(t reflect.Type) int {
	if s := _fixTypeSize(t); s > 0 {
		return s
	}
	var d Decoder = *this //avoid modify this
	switch t.Kind() {
	case reflect.String:
		var s __cntType = d.Uint32()
		size := int(s)
		return size + __cntSize //string length and data
	case reflect.Slice, reflect.Array:
		var s __cntType = d.Uint32()
		cnt := int(s)
		e := t.Elem()
		if s := _fixTypeSize(e); s > 0 {
			if t.Elem().Kind() == reflect.Bool { //compressed bool array
				return sizeofBoolArray(cnt)
			}
			return __cntSize + cnt*s
		} else {
			sum := __cntSize //array size
			for i, n := 0, cnt; i < n; i++ {
				s := d.sizeofType(e)
				d.Skip(s) //move to next element
				if s < 0 {
					return -1
				}
				sum += s
			}
			return sum
		}
	case reflect.Map:
		var s __cntType = d.Uint32()
		cnt := int(s)
		kt := t.Key()
		vt := t.Elem()
		sum := __cntSize //array size
		for i, n := 0, cnt; i < n; i++ {
			sk := d.sizeofType(kt)
			sv := d.sizeofType(vt)
			sum += (sk + sv)
		}
		return sum

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := d.sizeofType(t.Field(i).Type)
			d.Skip(s) //move to next element
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	}
	return -1
}

// decode bool array
func (this *Decoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			var _l __cntType = this.Uint32()
			l := int(_l)
			if k == reflect.Slice && l > 0 { //make a new slice
				v.Set(reflect.MakeSlice(v.Type(), l, l))
			}
			var b []byte
			for i := 0; i < l; i++ {
				_, bit := i/8, i%8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = this.reserve(1)
				}
				x := ((b[0] & mask) != 0)
				v.Index(i).SetBool(x)
			}
			return sizeofBoolArray(l)
		}
	}
	return -1
}
