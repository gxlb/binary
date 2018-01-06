package binary

import (
	"fmt"
	"io"
	"math"
	"reflect"
)

// NewDecoder make a new Decoder object with buffer.
func NewDecoder(buffer []byte) *Decoder {
	return NewDecoderEndian(buffer, DefaultEndian)
}

// NewDecoderEndian make a new Decoder object with buffer and endian.
func NewDecoderEndian(buffer []byte, endian Endian) *Decoder {
	p := &Decoder{}
	p.Init(buffer, endian)
	return p
}

// Decoder is used to decode byte array to go data.
type Decoder struct {
	coder
	reader    io.Reader //for decode from reader only
	boolValue byte      //last bool value byte
}

// Skip ignore the next size of bytes for encoding/decoding.
// It will panic If space not enough.
// It will return -1 if size <= 0.
func (decoder *Decoder) Skip(size int) int {
	if nil == decoder.mustReserve(size) {
		return -1
	}
	return size
}

// reserve returns next size bytes for encoding/decoding.
// It will panic if errors.
// do not use decoder.reserve, it was costly.
func (decoder *Decoder) mustReserve(size int) []byte {
	if decoder.reader != nil { //decode from reader
		if size > len(decoder.buff) {
			decoder.buff = make([]byte, size)
		}
		buff := decoder.buff[:size]
		if n, _ := decoder.reader.Read(buff); n < size {
			panic(io.ErrUnexpectedEOF)
		}
		return buff
	}
	return decoder.coder.mustReserve(size) //decode from bytes buffer
}

// reserve returns next size bytes for encoding/decoding.
func (decoder *Decoder) reserve(size int) ([]byte, error) {
	if decoder.reader != nil { //decode from reader
		if size > len(decoder.buff) {
			decoder.buff = make([]byte, size)
		}
		buff := decoder.buff[:size]
		if n, _ := decoder.reader.Read(buff); n < size {
			return nil, io.ErrUnexpectedEOF
		}
		return buff, nil
	}

	return decoder.coder.reserve(size) //decode from bytes buffer
}

// Init initialize Encoder with buffer and endian.
func (decoder *Decoder) Init(buffer []byte, endian Endian) {
	decoder.buff = buffer
	decoder.pos = 0
	decoder.endian = endian
}

// Bool decode a bool value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Bool() (value bool) {
	var err error
	if decoder.boolBit == 0 {
		var b []byte
		if b, err = decoder.reserve(1); err != nil {
			return false
		}
		decoder.boolValue = b[0]
	}

	mask := byte(1 << decoder.boolBit)
	decoder.boolBit = (decoder.boolBit + 1) % 8

	x := ((decoder.boolValue & mask) != 0)
	return x
}

// Int8 decode an int8 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Int8() int8 {
	return int8(decoder.Uint8())
}

// Uint8 decode a uint8 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Uint8() uint8 {
	b := decoder.mustReserve(1)
	x := b[0]
	return x
}

// Int16 decode an int16 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Int16(packed bool) int16 {
	if packed {
		x, _ := decoder.Varint()
		return int16(x)
	}

	return int16(decoder.Uint16(false))
}

// Uint16 decode a uint16 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Uint16(packed bool) uint16 {
	if packed {
		x, _ := decoder.Uvarint()
		return uint16(x)
	}

	b := decoder.mustReserve(2)
	x := decoder.endian.Uint16(b)
	return x
}

// Int32 decode an int32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Int32(packed bool) int32 {
	if packed {
		x, _ := decoder.Varint()
		return int32(x)
	}

	return int32(decoder.Uint32(false))
}

// Uint32 decode a uint32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Uint32(packed bool) uint32 {
	if packed {
		x, _ := decoder.Uvarint()
		return uint32(x)
	}

	b := decoder.mustReserve(4)
	x := decoder.endian.Uint32(b)
	return x
}

// Int64 decode an int64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Int64(packed bool) int64 {
	if packed {
		x, _ := decoder.Varint()
		return x
	}

	return int64(decoder.Uint64(false))
}

// Uint64 decode a uint64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Uint64(packed bool) uint64 {
	if packed {
		x, _ := decoder.Uvarint()
		return x
	}

	b := decoder.mustReserve(8)
	x := decoder.endian.Uint64(b)
	return x
}

// Float32 decode a float32 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Float32() float32 {
	x := math.Float32frombits(decoder.Uint32(false))
	return x
}

// Float64 decode a float64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Float64() float64 {
	x := math.Float64frombits(decoder.Uint64(false))
	return x
}

// Complex64 decode a complex64 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Complex64() complex64 {
	x := complex(decoder.Float32(), decoder.Float32())
	return x
}

// Complex128 decode a complex128 value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) Complex128() complex128 {
	x := complex(decoder.Float64(), decoder.Float64())
	return x
}

// String decode a string value from Decoder buffer.
// It will panic if buffer is not enough.
func (decoder *Decoder) String() string {
	s, _ := decoder.Uvarint()
	size := int(s)
	b := decoder.mustReserve(size)
	return string(b)
}

// Int decode an int value from Decoder buffer.
// It will panic if buffer is not enough.
// It use Varint() to decode as varint(1~10 bytes)
func (decoder *Decoder) Int() int {
	n, _ := decoder.Varint()
	return int(n)
}

// Uint decode a uint value from Decoder buffer.
// It will panic if buffer is not enough.
// It use Uvarint() to decode as uvarint(1~10 bytes)
func (decoder *Decoder) Uint() uint {
	n, _ := decoder.Uvarint()
	return uint(n)
}

// Varint decode an int64 value from Decoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
func (decoder *Decoder) Varint() (int64, int) {
	ux, n := decoder.Uvarint() // ok to continue in presence of error
	return ToVarint(ux), n
}

// Uvarint decode a uint64 value from Decoder buffer with varint(1~10 bytes).
// It will panic if buffer is not enough.
// It will return n <= 0 if varint error
func (decoder *Decoder) Uvarint() (uint64, int) {
	var x uint64
	var bit uint
	var i int
	for i = 0; i < MaxVarintLen64; i++ {
		b := decoder.Uint8()
		x |= uint64(b&0x7f) << bit
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				break // overflow
			}
			return x, i + 1
		}
		bit += 7
	}
	//return 0, 0
	panic(fmt.Errorf("binary.Decoder.Uvarint: overflow 64-bits value(pos:%d/%d)", decoder.Len(), decoder.Cap()))
}

func (decoder *Decoder) varint() (int64, int, error) {
	ux, n, err := decoder.uvarint() // ok to continue in presence of error
	return ToVarint(ux), n, err
}

func (decoder *Decoder) uvarint() (uint64, int, error) {
	var x uint64
	var bit uint
	var i int
	for i = 0; i < MaxVarintLen64; i++ {
		b := decoder.Uint8()
		x |= uint64(b&0x7f) << bit
		if b < 0x80 {
			if i >= MaxVarintLen64 || i == MaxVarintLen64-1 && b > 1 {
				break // overflow
			}
			return x, i + 1, nil
		}
		bit += 7
	}
	return 0, 0, fmt.Errorf("binary.Decoder.Uvarint: overflow 64-bits value(pos:%d/%d)", decoder.Len(), decoder.Cap())
}

// Value decode an interface value from Encoder buffer.
// x must be interface of pointer for modify.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (decoder *Decoder) Value(x interface{}) (err error) {
	return decoder.ValueX(x, defaultSerializer)
}

// ValueX decode an interface value from Encoder buffer.
// x must be interface of pointer for modify.
// enableSerializer switch if need check BinarySerilizer.
// It will return none-nil error if x contains unsupported types
// or buffer is not enough.
// It will check if x implements interface BinaryEncoder and use x.Encode first.
func (decoder *Decoder) ValueX(x interface{}, enableSerializer bool) (err error) {
	//	defer func() {
	//		if info := recover(); info != nil {
	//			err = info.(error)
	//			assert(err != nil, info)
	//		}
	//	}()

	//decoder.resetBoolCoder() //reset bool reader
	decoder.boolPos = -1
	decoder.boolBit = 0

	if decoder.fastValue(x) { //fast value path
		return nil
	}

	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr { //only support decode for pointer interface
		return decoder.value(v, 0, false, toplvSerializer(enableSerializer))
	}

	return typeError("binary.Decoder.Value: non-pointer type %s", v.Type(), true)
}

// use BinarySerializer interface to decode this value
func (decoder *Decoder) useSerializer(v reflect.Value) error {
	return decoder.Serializer(v.Interface())
}

// Serializer decode BinarySerializer x.
func (decoder *Decoder) Serializer(x interface{}) (err error) {
	//	t := reflect.TypeOf(x)
	//	if _, _, _, err := deepRegableType(t, true); err != nil {
	//		return err
	//	}
	if p, ok := x.(BinarySerializer); ok {
		size := p.Size()
		if err = p.Decode(decoder.buff[decoder.pos:]); err != nil {
			return err
		}
		_, err = decoder.reserve(size)
		return err
	}

	return typeError("binary: expect implements BinarySerializer %s", reflect.TypeOf(x), true)

}

// valuex decode v with serializer check
func (decoder *Decoder) value(v reflect.Value, depth int, packed bool, serializer serializerSwitch) error {
	k := v.Kind()
	if serializer.checkOk() ||
		serializer.needCheck() && k != reflect.Ptr && querySerializer(v.Type()) {
		return decoder.useSerializer(v.Addr())
	}

	switch k {
	case reflect.Int:
		v.SetInt(int64(decoder.Int()))
	case reflect.Uint:
		v.SetUint(uint64(decoder.Uint()))
	case reflect.Bool:
		v.SetBool(decoder.Bool())
	case reflect.Int8:
		v.SetInt(int64(decoder.Int8()))
	case reflect.Int16:
		v.SetInt(int64(decoder.Int16(packed)))
	case reflect.Int32:
		v.SetInt(int64(decoder.Int32(packed)))
	case reflect.Int64:
		v.SetInt(decoder.Int64(packed))
	case reflect.Uint8:
		v.SetUint(uint64(decoder.Uint8()))
	case reflect.Uint16:
		v.SetUint(uint64(decoder.Uint16(packed)))
	case reflect.Uint32:
		if packed {
			x, _ := decoder.Uvarint()
			v.SetUint(x)
		} else {
			b := decoder.mustReserve(4)
			x := decoder.endian.Uint32(b)
			v.SetUint(uint64(x))
		}
		//v.SetUint(uint64(decoder.Uint32(packed)))
	case reflect.Uint64:
		if packed {
			x, _ := decoder.Uvarint()
			v.SetUint(x)
		} else {
			b := decoder.mustReserve(8)
			x := decoder.endian.Uint64(b)
			v.SetUint(x)
		}

		//v.SetUint(decoder.Uint64(packed))
	case reflect.Float32:
		v.SetFloat(float64(decoder.Float32()))
	case reflect.Float64:
		v.SetFloat(decoder.Float64())
	case reflect.Complex64:
		v.SetComplex(complex128(decoder.Complex64()))
	case reflect.Complex128:
		v.SetComplex(decoder.Complex128())
	case reflect.String:
		v.SetString(decoder.String())

	case reflect.Slice, reflect.Array:
		elemT := v.Type().Elem()
		if !validUserType(elemT) { //verify array element is valid
			return fmt.Errorf("binary.Decoder.Value: unsupported type %s", v.Type().String())
		}

		elemSerializer := serializer.subSwitchCheck(elemT)
		if decoder.boolArray(v) < 0 { //deal with bool array first
			s, _ := decoder.Uvarint()
			size := int(s)
			if size > 0 && k == reflect.Slice { //make a new slice
				ns := reflect.MakeSlice(v.Type(), size, size)
				v.Set(ns)
			}

			for i, l := 0, v.Len(); i < size; i++ {
				if i < l {
					decoder.value(v.Index(i), depth+1, packed, elemSerializer) //ignore error
				} else {
					skiped := decoder.skipByType(v.Type().Elem(), packed, elemSerializer)
					//I'm sure here cannot find unsupported type
					assert(skiped >= 0, v.Type().Elem().String())
				}
			}
		}
	case reflect.Map:
		t := v.Type()
		kt, vt := t.Key(), t.Elem()
		//verify map key and value type are both valid
		if !validUserType(kt) || !validUserType(vt) {
			return typeError("binary.Decoder.Value: unsupported type %s", v.Type(), true)
		}

		if v.IsNil() {
			newmap := reflect.MakeMap(v.Type())
			v.Set(newmap)
		}

		keySerilaizer := serializer.subSwitchCheck(kt)
		valueSerilaizer := serializer.subSwitchCheck(vt)
		s, _ := decoder.Uvarint()
		for i, size := 0, int(s); i < size; i++ {
			key := reflect.New(kt).Elem()
			value := reflect.New(vt).Elem()
			decoder.value(key, depth+1, packed, keySerilaizer)     //ignore error
			decoder.value(value, depth+1, packed, valueSerilaizer) //ignore error
			v.SetMapIndex(key, value)
		}

	case reflect.Struct:
		return queryStruct(v.Type()).decode(decoder, v, serializer)

	default:
		if newPtr(v, decoder, depth) {
			if !v.IsNil() {
				return decoder.value(v.Elem(), depth+1, packed, serializer)
			}
		} else {
			return typeError("binary.Decoder.Value: unsupported type %s", v.Type(), true)
		}
	}
	return nil
}

func (decoder *Decoder) fastValue(x interface{}) bool {
	switch d := x.(type) {
	case *int:
		ux, _ := decoder.Uvarint() // ok to continue in presence of error
		*d = int(ToVarint(ux))
		//*d = decoder.Int()
	case *uint:
		ux, _ := decoder.Uvarint() // ok to continue in presence of error
		*d = uint(ux)
		//*d = decoder.Uint()
	case *bool:
		*d = decoder.Bool()
	case *int8:
		*d = decoder.Int8()
	case *uint8:
		*d = decoder.Uint8()
	case *int16:
		*d = decoder.Int16(false)
	case *uint16:
		*d = decoder.Uint16(false)
	case *int32:
		*d = decoder.Int32(false)
	case *uint32:
		*d = decoder.Uint32(false)
	case *float32:
		*d = decoder.Float32()
	case *int64:
		*d = decoder.Int64(false)
	case *uint64:
		*d = decoder.Uint64(false)
	case *float64:
		*d = decoder.Float64()
	case *complex64:
		*d = decoder.Complex64()
	case *complex128:
		*d = decoder.Complex128()
	case *string:
		*d = decoder.String()
	case *[]bool:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]bool, l)
		var b []byte
		for i := 0; i < l; i++ {
			_, bit := i/8, i%8
			mask := byte(1 << uint(bit))
			if bit == 0 {
				b = decoder.mustReserve(1)
			}
			x := ((b[0] & mask) != 0)
			(*d)[i] = x
		}
	case *[]int:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]int, l)
		for i := 0; i < l; i++ {
			ux, _ := decoder.Uvarint() // ok to continue in presence of error
			(*d)[i] = int(ToVarint(ux))
			//(*d)[i] = decoder.Int()
		}
	case *[]uint:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]uint, l)
		for i := 0; i < l; i++ {
			ux, _ := decoder.Uvarint() // ok to continue in presence of error
			(*d)[i] = uint(ux)
			//(*d)[i] = decoder.Uint()
		}

	case *[]int8:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]int8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Int8()
		}
	case *[]uint8:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]uint8, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Uint8()
		}
	case *[]int16:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]int16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Int16(false)
		}
	case *[]uint16:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]uint16, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Uint16(false)
		}
	case *[]int32:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]int32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Int32(false)
		}
	case *[]uint32:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]uint32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Uint32(false)
		}
	case *[]int64:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]int64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Int64(false)
		}
	case *[]uint64:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]uint64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Uint64(false)
		}
	case *[]float32:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]float32, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Float32()
		}
	case *[]float64:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]float64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Float64()
		}
	case *[]complex64:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]complex64, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Complex64()
		}
	case *[]complex128:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]complex128, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.Complex128()
		}
	case *[]string:
		s, _ := decoder.Uvarint()
		l := int(s)
		*d = make([]string, l)
		for i := 0; i < l; i++ {
			(*d)[i] = decoder.String()
		}
	default:
		return false
	}
	return true
}

//TODO:
// serializer bug
func (decoder *Decoder) skipByType(t reflect.Type, packed bool, serializer serializerSwitch) int {
	if s := fixedTypeSize(t); s > 0 {
		if packedType := packedIntsType(t); packedType > 0 && packed {
			switch packedType {
			case _SignedInts:
				_, n := decoder.Varint()
				return n
			case _UnsignedInts:
				_, n := decoder.Uvarint()
				return n
			}
		} else {
			decoder.Skip(s)
			return s
		}
	}
	switch t.Kind() {
	case reflect.Ptr:
		if isNotNil := decoder.Bool(); isNotNil {
			return decoder.skipByType(t.Elem(), packed, serializer) + 1
		}
		return 1
	case reflect.Bool:
		decoder.Bool()
		return 1
	case reflect.Int:
		_, n := decoder.Varint()
		return n
	case reflect.Uint:
		_, n := decoder.Uvarint()
		return n
	case reflect.String:
		s, n := decoder.Uvarint()
		size := int(s) //string length and data
		decoder.Skip(size)
		return size + n
	case reflect.Slice, reflect.Array:
		s, sLen := decoder.Uvarint()
		cnt := int(s)
		elemtype := t.Elem()
		elemSerializer := serializer.subSwitchCheck(elemtype)
		if s := fixedTypeSize(elemtype); s > 0 {
			size := cnt * s
			decoder.Skip(size)
			return size
		}

		if elemtype.Kind() == reflect.Bool { //compressed bool array
			totalSize := sizeofBoolArray(cnt)
			size := totalSize - SizeofUvarint(uint64(cnt)) //cnt has been read
			decoder.Skip(size)
			return totalSize
		}

		sum := sLen //array size
		for i, n := 0, cnt; i < n; i++ {
			s := decoder.skipByType(elemtype, packed, elemSerializer)
			assert(s >= 0, "skip fail: "+elemtype.String()) //I'm sure here cannot find unsupported type
			sum += s
		}
		return sum
	case reflect.Map:
		s, sLen := decoder.Uvarint()
		cnt := int(s)
		kt := t.Key()
		vt := t.Elem()
		sum := sLen //array size
		keySerilaizer := serializer.subSwitchCheck(kt)
		valueSerilaizer := serializer.subSwitchCheck(vt)
		for i, n := 0, cnt; i < n; i++ {
			sum += decoder.skipByType(kt, packed, keySerilaizer)
			sum += decoder.skipByType(vt, packed, valueSerilaizer)
		}
		return sum

	case reflect.Struct:
		return queryStruct(t).decodeSkipByType(decoder, t, packed, serializer)
	}
	return -1
}

// decode bool array
func (decoder *Decoder) boolArray(v reflect.Value) int {
	if k := v.Kind(); k == reflect.Slice || k == reflect.Array {
		if v.Type().Elem().Kind() == reflect.Bool {
			_l, _ := decoder.Uvarint()
			l := int(_l)
			if k == reflect.Slice && l > 0 { //make a new slice
				v.Set(reflect.MakeSlice(v.Type(), l, l))
			}
			var b []byte
			for i := 0; i < l; i++ {
				_, bit := i/8, i%8
				mask := byte(1 << uint(bit))
				if bit == 0 {
					b = decoder.mustReserve(1)
				}
				x := ((b[0] & mask) != 0)
				v.Index(i).SetBool(x)
			}
			return sizeofBoolArray(l)
		}
	}
	return -1
}
