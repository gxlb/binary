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

// NewEncoderBuffer make a new Encoder object with buffer.
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
func (encoder *Encoder) Init(size int, endian Endian) {
	encoder.buff = make([]byte, size)
	encoder.pos = 0
	encoder.endian = endian
}

// ResizeBuffer confirm that len(buffer) >= size and alloc larger buffer if necessary
// It will call Reset to initial encoder state of buffer
func (encoder *Encoder) ResizeBuffer(size int) bool {
	ok := len(encoder.buff) < size
	if ok {
		encoder.buff = make([]byte, size)
	}
	encoder.Reset()
	return ok
}

// Bool encode a bool value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Bool(x bool) {
	if encoder.boolBit == 0 {
		b := encoder.reserve(1)
		b[0] = 0
		encoder.boolPos = encoder.pos - 1
	}

	if mask := byte(1 << encoder.boolBit); x {
		encoder.buff[encoder.boolPos] |= mask
	}
	encoder.boolBit = (encoder.boolBit + 1) % 8
}

// Int8 encode an int8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int8(x int8) {
	encoder.Uint8(uint8(x))
}

// Uint8 encode a uint8 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint8(x uint8) {
	b := encoder.reserve(1)
	b[0] = x
}

// Int16 encode an int16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int16(x int16, packed bool) {
	if packed {
		encoder.Varint(int64(x))
	} else {
		encoder.Uint16(uint16(x), false)
	}

}

// Uint16 encode a uint16 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint16(x uint16, packed bool) {
	if packed {
		encoder.Uvarint(uint64(x))
	} else {
		b := encoder.reserve(2)
		encoder.endian.PutUint16(b, x)
	}
}

// Int32 encode an int32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int32(x int32, packed bool) {
	if packed {
		encoder.Varint(int64(x))
	} else {
		encoder.Uint32(uint32(x), false)
	}
}

// Uint32 encode a uint32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint32(x uint32, packed bool) {
	if packed {
		encoder.Uvarint(uint64(x))
	} else {
		b := encoder.reserve(4)
		encoder.endian.PutUint32(b, x)
	}
}

// Int64 encode an int64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Int64(x int64, packed bool) {
	if packed {
		encoder.Varint(x)
	} else {
		encoder.Uint64(uint64(x), false)
	}
}

// Uint64 encode a uint64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Uint64(x uint64, packed bool) {
	if packed {
		encoder.Uvarint(x)
	} else {
		b := encoder.reserve(8)
		encoder.endian.PutUint64(b, x)
	}
}

// Float32 encode a float32 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Float32(x float32) {
	encoder.Uint32(math.Float32bits(x), false)
}

// Float64 encode a float64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Float64(x float64) {
	encoder.Uint64(math.Float64bits(x), false)
}

// Complex64 encode a complex64 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Complex64(x complex64) {
	encoder.Uint32(math.Float32bits(real(x)), false)
	encoder.Uint32(math.Float32bits(imag(x)), false)
}

// Complex128 encode a complex128 value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) Complex128(x complex128) {
	encoder.Uint64(math.Float64bits(real(x)), false)
	encoder.Uint64(math.Float64bits(imag(x)), false)
}

// String encode a string value to Encoder buffer.
// It will panic if buffer is not enough.
func (encoder *Encoder) String(x string) {
	_b := []byte(x)
	size := len(_b)
	encoder.Uvarint(uint64(size))
	buff := encoder.reserve(size)
	copy(buff, _b)
}

// Int encode an int value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Varint() to encode as varint(1~10 bytes)
func (encoder *Encoder) Int(x int) {
	encoder.Varint(int64(x))
}

// Uint encode a uint value to Encoder buffer.
// It will panic if buffer is not enough.
// It use Uvarint() to encode as uvarint(1~10 bytes)
func (encoder *Encoder) Uint(x uint) {
	encoder.Uvarint(uint64(x))
}

// Varint encode an int64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (encoder *Encoder) Varint(x int64) int {
	return encoder.Uvarint(ToUvarint(x))
}

// Uvarint encode a uint64 value to Encoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (encoder *Encoder) Uvarint(x uint64) int {
	i, _x := 0, x
	for ; _x >= 0x80; _x >>= 7 {
		encoder.Uint8(byte(_x) | 0x80)
		i++
	}
	encoder.Uint8(byte(_x))
	return i + 1
}

// Value encode an interface value to Encoder buffer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (encoder *Encoder) Value(x interface{}) (err error) {
	return encoder.ValueX(x, defaultSerializer)
}

// ValueX encode an interface value to Encoder buffer.
// enableSerializer switch if need check BinarySerilizer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (encoder *Encoder) ValueX(x interface{}, enableSerializer bool) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	encoder.resetBoolCoder() //reset bool writer

	if encoder.fastValue(x) { //fast value path
		return nil
	}

	v := reflect.ValueOf(x)

	//	if p, ok := x.(BinaryEncoder); ok {
	//		if _, _ok := x.(BinarySizer); !_ok { //interface verification
	//			panic(fmt.Errorf("expect but not BinarySizer: %s", v.Type().String()))
	//		}

	//		r, err := p.Encode(encoder.buff[encoder.pos:])
	//		if err == nil {
	//			encoder.reserve(len(r))
	//		}
	//		return err
	//	}

	//	if _, _ok := x.(BinarySizer); _ok { //interface verification
	//		panic(fmt.Errorf("unexpected BinarySizer: %s", v.Type().String()))
	//	}

	return encoder.value(reflect.Indirect(v), false, toplvSerializer(enableSerializer))
}

func (encoder *Encoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case int:
		encoder.Int(d)
	case uint:
		encoder.Uint(d)

	case bool:
		encoder.Bool(d)
	case int8:
		encoder.Int8(d)
	case uint8:
		encoder.Uint8(d)
	case int16:
		encoder.Int16(d, false)
	case uint16:
		encoder.Uint16(d, false)
	case int32:
		encoder.Int32(d, false)
	case uint32:
		encoder.Uint32(d, false)
	case float32:
		encoder.Float32(d)
	case int64:
		encoder.Int64(d, false)
	case uint64:
		encoder.Uint64(d, false)
	case float64:
		encoder.Float64(d)
	case complex64:
		encoder.Complex64(d)
	case complex128:
		encoder.Complex128(d)
	case string:
		encoder.String(d)
	case []bool:
		l := len(d)
		encoder.Uvarint(uint64(l))
		var b []byte
		for i := 0; i < l; i++ {
			bit := i % 8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b = encoder.reserve(1)
				b[0] = 0
			}
			if x := d[i]; x {
				b[0] |= mask
			}
		}

	case []int8:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int8(d[i])
		}
	case []uint8:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint8(d[i])
		}
	case []int16:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int16(d[i], false)
		}
	case []uint16:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint16(d[i], false)
		}

	case []int32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int32(d[i], false)
		}
	case []uint32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint32(d[i], false)
		}
	case []int64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int64(d[i], false)
		}
	case []uint64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint64(d[i], false)
		}
	case []float32:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Float32(d[i])
		}
	case []float64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Float64(d[i])
		}
	case []complex64:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Complex64(d[i])
		}
	case []complex128:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Complex128(d[i])
		}
	case []string:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.String(d[i])
		}
	case []int:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Int(d[i])
		}
	case []uint:
		l := len(d)
		encoder.Uvarint(uint64(len(d)))
		for i := 0; i < l; i++ {
			encoder.Uint(d[i])
		}
	default:
		return false
	}
	return true

}

// use BinarySerializer interface to encode this value
func (encoder *Encoder) useSerializer(v reflect.Value) error {
	x := v.Interface()
	if p, ok := x.(BinaryEncoder); ok {
		r, err := p.Encode(encoder.buff[encoder.pos:])
		if err == nil {
			encoder.reserve(len(r))
		}
		return err
	}

	panic(typeError("expect BinarySerializer %s", v.Type(), true))
}

func (encoder *Encoder) value(v reflect.Value, packed bool, serializer SerializerSwitch) error {
	// check Packer interface for every value is perfect
	// but encoder is too costly
	//
	//	if t := v.Type(); t.Implements(tPacker) {
	//		if !t.Implements(tSizer) { //interface verification
	//			panic(fmt.Errorf("pected but not Sizer: %s", t.String()))
	//		}
	//		packer := v.Interface().(Packer)
	//		reault, err := packer.Pack(encoder.buff[encoder.pos:])
	//		if err == nil {
	//			encoder.reserve(len(reault))
	//		}
	//		return err
	//	} else {
	//		if t.Implements(tSizer) { //interface verification
	//			panic(fmt.Errorf("unexpected Sizer: %s", v.Type().String()))
	//		}
	//	}

	k := v.Kind()
	if serializer.CheckOk() ||
		serializer.NeedCheck() && k != reflect.Ptr && querySerializer(v.Type()) {
		return encoder.useSerializer(v)
	}

	switch k {
	case reflect.Int:
		encoder.Int(int(v.Int()))
	case reflect.Uint:
		encoder.Uint(uint(v.Uint()))

	case reflect.Bool:
		encoder.Bool(v.Bool())

	case reflect.Int8:
		encoder.Int8(int8(v.Int()))
	case reflect.Int16:
		encoder.Int16(int16(v.Int()), packed)
	case reflect.Int32:
		encoder.Int32(int32(v.Int()), packed)
	case reflect.Int64:
		encoder.Int64(v.Int(), packed)

	case reflect.Uint8:
		encoder.Uint8(uint8(v.Uint()))
	case reflect.Uint16:
		encoder.Uint16(uint16(v.Uint()), packed)
	case reflect.Uint32:
		encoder.Uint32(uint32(v.Uint()), packed)
	case reflect.Uint64:
		encoder.Uint64(v.Uint(), packed)

	case reflect.Float32:
		encoder.Float32(float32(v.Float()))
	case reflect.Float64:
		encoder.Float64(v.Float())

	case reflect.Complex64:
		x := v.Complex()
		encoder.Complex64(complex64(x))

	case reflect.Complex128:
		x := v.Complex()
		encoder.Complex128(x)

	case reflect.String:
		encoder.String(v.String())

	case reflect.Slice, reflect.Array:
		elemT := v.Type().Elem()
		if !validUserType(elemT) { //verify array element is valid
			return fmt.Errorf("binary.Encoder.Value: unsupported type %s", v.Type().String())
		}
		elemSerializer := serializer.SubSwitchCheck(elemT)
		if encoder.boolArray(v) < 0 { //deal with bool array first
			l := v.Len()
			encoder.Uvarint(uint64(l))
			for i := 0; i < l; i++ {
				assert(encoder.value(v.Index(i), packed, elemSerializer) == nil, "")
			}
		}
	case reflect.Map:
		t := v.Type()
		kt := t.Key()
		vt := t.Elem()
		if !validUserType(kt) || !validUserType(vt) { //verify map key and value type are both valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}

		keySerilaizer := serializer.SubSwitchCheck(kt)
		valueSerilaizer := serializer.SubSwitchCheck(vt)

		keys := v.MapKeys()
		l := len(keys)
		encoder.Uvarint(uint64(l))
		for i := 0; i < l; i++ {
			key := keys[i]
			assert(encoder.value(key, packed, keySerilaizer) == nil, "")
			assert(encoder.value(v.MapIndex(key), packed, valueSerilaizer) == nil, "")
		}
	case reflect.Struct:
		return queryStruct(v.Type()).encode(encoder, v, serializer)

	case reflect.Ptr:
		if !validUserType(v.Type()) {
			return fmt.Errorf("binary.Encoder.Value: unsupported type %s", v.Type().String())
		}
		if !v.IsNil() {
			encoder.Bool(true)
			if e := v.Elem(); e.Kind() != reflect.Ptr {
				return encoder.value(e, packed, serializer)
			}
		} else {
			encoder.Bool(false)
			//			if encoder.nilPointer(v.Type()) < 0 {
			//				return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Type().String())
			//			}
		}
		//	case reflect.Invalid://BUG: it will panic to get zero.Type
		//		return fmt.Errorf("binary.Encoder.Value: unsupported type [%s]", v.Kind().String())
	default:
		return typeError("binary.Encoder.Value: unsupported type [%s]", v.Type(), true)
	}
	return nil
}

// encode bool array
func (encoder *Encoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			l := v.Len()
			encoder.Uvarint(uint64(l))
			var b []byte
			for i := 0; i < l; i++ {
				bit := i % 8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = encoder.reserve(1)
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
