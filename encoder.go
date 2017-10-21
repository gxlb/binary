package binary

import (
	"fmt"
	"math"
	"reflect"
)

// NewEncoder make a new Encoder object with buffer size.
func NewEncoder(size int) *Encoder {
	return NewEncoderEndian(size, DefaultEndian)
}

// NewEncoder make a new Encoder object with buffer.
func NewEncoderBuffer(buffer []byte) *Encoder {
	p := &Encoder{}
	//assert(buffer != nil, "nil buffer")
	p.buff = buffer
	p.endian = DefaultEndian
	p.pos = 0
	return p
}

// NewEncoderEndian make a new Encoder object with buffer size and endian.
func NewEncoderEndian(size int, endian Endian) *Encoder {
	p := &Encoder{}
	p.Init(size, endian)
	return p
}

// Encoder is used to encode go data to byte array.
type Encoder struct {
	coder
}

// Init initialize Encoder with buffer size and endian.
func (this *Encoder) Init(size int, endian Endian) {
	this.buff = make([]byte, size)
	this.pos = 0
	this.endian = endian
}

// Bool encode a bool value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Bool(x bool) {
	b := this.reserve(1)
	if x {
		b[0] = 1
	} else {
		b[0] = 0
	}
}

// Int8 encode an int8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Int8(x int8) {
	this.Uint8(uint8(x))
}

// Uint8 encode a uint8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Uint8(x uint8) {
	b := this.reserve(1)
	b[0] = x
}

// Int16 encode an int16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Int16(x int16) {
	this.Uint16(uint16(x))
}

// Uint16 encode a uint16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Uint16(x uint16) {
	b := this.reserve(2)
	this.endian.PutUint16(b, x)
}

// Int32 encode an int32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Int32(x int32) {
	this.Uint32(uint32(x))
}

// Uint32 encode a uint32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Uint32(x uint32) {
	b := this.reserve(4)
	this.endian.PutUint32(b, x)
}

// Int64 encode an int64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Int64(x int64) {
	this.Uint64(uint64(x))
}

// Uint64 encode a uint64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Uint64(x uint64) {
	b := this.reserve(8)
	this.endian.PutUint64(b, x)
}

// Float32 encode a float32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Float32(x float32) {
	this.Uint32(math.Float32bits(x))
}

// Float64 encode a float64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Float64(x float64) {
	this.Uint64(math.Float64bits(x))
}

// Complex64 encode a complex64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Complex64(x complex64) {
	this.Uint32(math.Float32bits(real(x)))
	this.Uint32(math.Float32bits(imag(x)))
}

// Complex128 encode a complex128 value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) Complex128(x complex128) {
	this.Uint64(math.Float64bits(real(x)))
	this.Uint64(math.Float64bits(imag(x)))
}

// String encode a string value to Encoder buffer.
// It will panic if buffer is not enough.
func (this *Encoder) String(x string) {
	_b := []byte(x)
	size := len(_b)
	this.Uvarint(uint64(size))
	buff := this.reserve(size)
	copy(buff, _b)
}

// Int encode an int value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Varint() to encode as varint(1~10 bytes)
func (this *Encoder) Int(x int) {
	this.Varint(int64(x))
}

// Uint encode a uint value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Uvarint() to encode as uvarint(1~10 bytes)
func (this *Encoder) Uint(x uint) {
	this.Uvarint(uint64(x))
}

// Varint encode an int64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (this *Encoder) Varint(x int64) int {
	return this.Uvarint(ToUvarint(x))
}

// Uvarint encode a uint64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (this *Encoder) Uvarint(x uint64) int {
	i, _x := 0, x
	for ; _x >= 0x80; _x >>= 7 {
		this.Uint8(byte(_x) | 0x80)
		i++
	}
	this.Uint8(byte(_x))
	return i + 1
}

// Value encode an interface value to Encoder buffer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
func (this *Encoder) Value(x interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	if this.fastValue(x) { //fast value path
		return nil
	}

	v := reflect.ValueOf(x)

	if p, ok := x.(Packer); ok {
		if _, _ok := x.(Sizer); !_ok { //interface verification
			panic(fmt.Errorf("pected but not Sizer: %s", v.Type().String()))
		}

		r, err := p.Pack(this.buff[this.pos:])
		if err == nil {
			this.reserve(len(r))
		}
		return err
	} else {
		if _, _ok := x.(Sizer); _ok { //interface verification
			panic(fmt.Errorf("unexpected Sizer: %s", v.Type().String()))
		}
	}

	return this.value(v)
}

func (this *Encoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case int:
		this.Int(d)
	case uint:
		this.Uint(d)

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
		this.Uvarint(uint64(l))
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
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Int8(d[i])
		}
	case []uint8:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Uint8(d[i])
		}
	case []int16:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Int16(d[i])
		}
	case []uint16:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Uint16(d[i])
		}

	case []int32:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Int32(d[i])
		}
	case []uint32:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Uint32(d[i])
		}
	case []int64:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Int64(d[i])
		}
	case []uint64:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Uint64(d[i])
		}
	case []float32:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Float32(d[i])
		}
	case []float64:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Float64(d[i])
		}
	case []complex64:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Complex64(d[i])
		}
	case []complex128:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Complex128(d[i])
		}
	case []string:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.String(d[i])
		}
	case []int:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Int(d[i])
		}
	case []uint:
		l := len(d)
		this.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			this.Uint(d[i])
		}
	default:
		return false
	}
	return true

}

func (this *Encoder) value(v reflect.Value) error {
	// check Packer interface for every value is perfect
	// but this is too costly
	//
	//	if t := v.Type(); t.Implements(tPacker) {
	//		if !t.Implements(tSizer) { //interface verification
	//			panic(fmt.Errorf("pected but not Sizer: %s", t.String()))
	//		}
	//		packer := v.Interface().(Packer)
	//		reault, err := packer.Pack(this.buff[this.pos:])
	//		if err == nil {
	//			this.reserve(len(reault))
	//		}
	//		return err
	//	} else {
	//		if t.Implements(tSizer) { //interface verification
	//			panic(fmt.Errorf("unexpected Sizer: %s", v.Type().String()))
	//		}
	//	}

	switch k := v.Kind(); k {
	case reflect.Int:
		this.Int(int(v.Int()))
	case reflect.Uint:
		this.Uint(uint(v.Uint()))

	case reflect.Bool:
		this.Bool(v.Bool())

	case reflect.Int8:
		this.Int8(int8(v.Int()))
	case reflect.Int16:
		this.Int16(int16(v.Int()))
	case reflect.Int32:
		this.Int32(int32(v.Int()))
	case reflect.Int64:
		this.Int64(v.Int())

	case reflect.Uint8:
		this.Uint8(uint8(v.Uint()))
	case reflect.Uint16:
		this.Uint16(uint16(v.Uint()))
	case reflect.Uint32:
		this.Uint32(uint32(v.Uint()))
	case reflect.Uint64:
		this.Uint64(v.Uint())

	case reflect.Float32:
		this.Float32(float32(v.Float()))
	case reflect.Float64:
		this.Float64(v.Float())

	case reflect.Complex64:
		x := v.Complex()
		this.Complex64(complex64(x))

	case reflect.Complex128:
		x := v.Complex()
		this.Complex128(x)

	case reflect.String:
		this.String(v.String())

	case reflect.Slice, reflect.Array:
		if sizeofNilPointer(v.Type().Elem()) < 0 { //verify array element is valid
			return fmt.Errorf("binary.Encoder.Value: unsupported type %s", v.Type().String())
		}
		if this.boolArray(v) < 0 { //deal with bool array first
			l := v.Len()
			this.Uvarint(uint64(l))
			for i := 0; i < l; i++ {
				this.value(v.Index(i))
			}
		}
	case reflect.Map:
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()
		if sizeofNilPointer(kt) < 0 ||
			sizeofNilPointer(vt) < 0 { //verify map key and value type are both valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}

		keys := v.MapKeys()
		l := len(keys)
		this.Uvarint(uint64(l))
		for i := 0; i < l; i++ {
			key := keys[i]
			this.value(key)
			this.value(v.MapIndex(key))
		}
	case reflect.Struct:
		return queryStruct(v.Type()).encode(this, v)

	case reflect.Ptr:
		if !v.IsNil() {
			if e := v.Elem(); e.Kind() != reflect.Ptr {
				return this.value(e)
			}
		} else {
			if this.nilPointer(v.Type()) < 0 {
				return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Type().String())
			}
		}
		//	case reflect.Invalid://BUG: it will panic to get zero.Type
		//		return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Kind().String())
	default:
		return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Type().String())
	}
	return nil
}

func (this *Encoder) nilPointer(t reflect.Type) int {
	tt := t
	if tt.Kind() == reflect.Ptr {
		tt = t.Elem()
		if tt.Kind() == reflect.Ptr {
			return -1
		}
	}
	if s := _fixTypeSize(tt); s > 0 { //fix size
		return this.Skip(s)
	}
	switch tt.Kind() {
	case reflect.Int, reflect.Uint: //zero varint will be encoded as 1 byte
		return this.Uvarint(0)
	case reflect.Slice, reflect.String:
		return this.Uvarint(0)
	case reflect.Array:
		l := tt.Len()
		n := this.Uvarint(uint64(l))
		if tt.Elem().Kind() == reflect.Bool { //bool array
			n2 := sizeofBoolArray(n)
			this.Skip(n2 - n)
			n = n2
		} else {
			tte := tt.Elem()
			for i := 0; i < l; i++ {
				n += this.nilPointer(tte)
			}
		}
		return n

	case reflect.Struct:
		return queryStruct(tt).encodeNilPointer(this, tt)
	}
	return -1
}

// encode bool array
func (this *Encoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			l := v.Len()
			this.Uvarint(uint64(l))
			var b []byte
			for i := 0; i < l; i++ {
				bit := i % 8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = this.reserve(1)
					b[0] = 0
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
