package binary

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func sizeof(data interface{}) int {
	if s := fastSizeof(data); s >= 0 {
		return s
	}

	s := bitsOfValue(reflect.ValueOf(data), true, false)
	if s < 0 {
		return -1
	}
	return (s + 7) / 8
}

func fastSizeof(data interface{}) int {
	switch d := data.(type) { //fast size calculation
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case int16, uint16, *int16, *uint16:
		return 2
	case int32, uint32, *int32, *uint32, float32, *float32:
		return 4
	case int64, uint64, *int64, *uint64, float64, *float64, complex64, *complex64:
		return 8
	case complex128, *complex128:
		return 16
	case string:
		return sizeofString(len(d))

	case int:
		return SizeofVarint(int64(d))
	case uint:
		return SizeofUvarint(uint64(d))
	case *int:
		if d != nil {
			return SizeofVarint(int64(*d))
		}
	case *uint:
		if d != nil {
			return SizeofUvarint(uint64(*d))
		}
	case []bool:
		return sizeofBoolArray(len(d))
	case []int8:
		return sizeofFixArray(len(d), 1)
	case []uint8:
		return sizeofFixArray(len(d), 1)
	case []int16:
		return sizeofFixArray(len(d), 2)
	case []uint16:
		return sizeofFixArray(len(d), 2)
	case []int32:
		return sizeofFixArray(len(d), 4)
	case []uint32:
		return sizeofFixArray(len(d), 4)
	case []float32:
		return sizeofFixArray(len(d), 4)
	case []int64:
		return sizeofFixArray(len(d), 8)
	case []uint64:
		return sizeofFixArray(len(d), 8)
	case []float64:
		return sizeofFixArray(len(d), 8)
	case []complex64:
		return sizeofFixArray(len(d), 8)
	case []complex128:
		return sizeofFixArray(len(d), 16)
	case []string:
		l := len(d)
		s := SizeofUvarint(uint64(l))
		for _, v := range d {
			s += sizeofString(len(v))
		}
		return s
	case []int:
		l := len(d)
		s := SizeofUvarint(uint64(l))
		for _, v := range d {
			s += SizeofVarint(int64(v))
		}
		return s
	case []uint:
		l := len(d)
		s := SizeofUvarint(uint64(l))
		for _, v := range d {
			s += SizeofUvarint(uint64(v))
		}
		return s

	case *[]bool:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]int8:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]uint8:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]int16:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]uint16:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]int32:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]uint32:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]float32:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]int64:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]uint64:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]float64:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]complex64:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]complex128:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]string:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]int:
		if d != nil {
			return fastSizeof(*d)
		}
	case *[]uint:
		if d != nil {
			return fastSizeof(*d)
		}
	}
	return -1
}

func assert(b bool, msg interface{}) {
	if !b {
		panic(fmt.Errorf("%s", msg))
	}
}

func bitsOfUnfixedArray(v reflect.Value, packed bool) int {
	if !validUserType(v.Type().Elem()) { //check if array element type valid
		return -1
	}

	arrayLen := v.Len()
	sum := SizeofUvarint(uint64(arrayLen)) * 8 //array size bytes num
	for i, n := 0, arrayLen; i < n; i++ {
		s := bitsOfValue(v.Index(i), false, packed)
		//assert(s >= 0, v.Type().String()) //element size must not error
		sum += s
	}
	return sum
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func bitsOfValue(v reflect.Value, topLevel bool, packed bool) (r int) {
	//	defer func() {
	//		fmt.Printf("bitsOfValue(%#v)=%d\n", v.Interface(), r)
	//	}()
	bits := 0
	if v.Kind() == reflect.Ptr { //nil is not aviable
		if !topLevel {
			bits = 1
		}
		if v.IsNil() {
			if topLevel || !validUserType(v.Type()) {
				return -1
			}
			return 1
		}
	}

	v = reflect.Indirect(v) //redrect pointer to it's value
	t := v.Type()
	if s := fixedTypeSize(t); s > 0 { //fixed size
		if packedType := packedIntsType(t); packedType > 0 && packed {
			switch packedType {
			case _SignedInts:
				return SizeofVarint(v.Int())*8 + bits
			case _UnsignedInts:
				return SizeofUvarint(v.Uint())*8 + bits
			}
		} else {
			return s*8 + bits
		}
	}
	switch t := v.Type(); t.Kind() {
	case reflect.Bool:
		return 1 + bits
	case reflect.Int:
		return SizeofVarint(v.Int())*8 + bits
	case reflect.Uint:
		return SizeofUvarint(v.Uint())*8 + bits
	case reflect.Slice, reflect.Array:
		arrayLen := v.Len()
		elemtype := t.Elem()
		if s := fixedTypeSize(elemtype); s > 0 {
			if packedIntsType(elemtype) > 0 && packed {
				return bitsOfUnfixedArray(v, packed) + bits
			}

			return sizeofFixArray(arrayLen, s)*8 + bits
		}

		if elemtype.Kind() == reflect.Bool {
			return sizeofBoolArray(arrayLen)*8 + bits
		}
		return bitsOfUnfixedArray(v, packed) + bits
	case reflect.Map:
		mapLen := v.Len()
		sum := SizeofUvarint(uint64(mapLen))*8 + bits //array size
		keys := v.MapKeys()

		if !validUserType(t.Key()) ||
			!validUserType(t.Elem()) { //check if map key and value type valid
			return -1
		}

		for i := 0; i < mapLen; i++ {
			key := keys[i]
			sizeKey := bitsOfValue(key, false, packed)
			//assert(sizeKey >= 0, key.Type().Kind().String()) //key size must not error

			sum += sizeKey
			value := v.MapIndex(key)
			sizeValue := bitsOfValue(value, false, packed)
			//assert(sizeValue >= 0, value.Type().Kind().String()) //key size must not error

			sum += sizeValue
		}
		return sum

	case reflect.Struct:
		return queryStruct(v.Type()).bitsOfValue(v) + bits

	case reflect.String:
		return sizeofString(v.Len())*8 + bits //string length and data
	}
	return -1
}

//BUG:
// bool as a byte, but not bit
// always use this functin to verify if Type is valid
// and do not care the value of return bytes
func sizeofNilPointer(t reflect.Type) int {
	tt := t
	if tt.Kind() == reflect.Ptr {
		tt = t.Elem()
	}
	if s := fixedTypeSize(tt); s > 0 { //fix size
		return s
	}
	switch tt.Kind() {
	case reflect.Bool:
		return 1
	case reflect.Int, reflect.Uint: //zero varint will be encoded as 1 byte
		return 1
	case reflect.String:
		return SizeofUvarint(0)
	case reflect.Slice:
		if validUserType(tt.Elem()) { //verify element type valid
			return SizeofUvarint(0)
		}
	case reflect.Map:
		if validUserType(tt.Key()) &&
			validUserType(tt.Elem()) { //verify key and value type valid
			return SizeofUvarint(0)
		}
	case reflect.Array:
		elemtype := tt.Elem()
		if s := fixedTypeSize(elemtype); s > 0 {
			return sizeofFixArray(tt.Len(), s)
		}

		if elemtype.Kind() == reflect.Bool {
			return sizeofBoolArray(tt.Len())
		}
		size := sizeofNilPointer(elemtype)
		if size > 0 { //verify element type valid
			return sizeofFixArray(tt.Len(), size)
		}
	case reflect.Struct:
		return queryStruct(tt).sizeofNilPointer(tt)
	}

	return -1
}

const (
	_SignedInts = iota + 1
	_UnsignedInts
)

func packedIntsType(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Int16, reflect.Int32, reflect.Int64:
		return _SignedInts
	case reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return _UnsignedInts
	}
	return 0
}

func fixedTypeSize(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64:
		return 8
	case reflect.Complex128:
		return 16
	}
	return -1
}

// Auto allocate for aviable pointer
func newPtr(v reflect.Value, decoder *Decoder, topLevel bool) bool {
	if v.Kind() == reflect.Ptr {
		e := v.Type().Elem()
		switch e.Kind() {
		case reflect.Array, reflect.Struct, reflect.Slice, reflect.Map:
			if !validUserType(e) { //check if valid pointer type
				return false
			}
			fallthrough
		case reflect.Int, reflect.Uint, reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16,
			reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64,
			reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64,
			reflect.Complex128, reflect.String:
			isNotNilPointer := false
			if !topLevel {
				isNotNilPointer = decoder.Bool()
				if v.IsNil() {
					if isNotNilPointer {
						v.Set(reflect.New(e))
					}
				}
			}
			return true
		}
	}
	return false
}

// NOTE:
// This function will make the encode/decode of struct slow down.
// It is recommended to use RegStruct to improve this case.
func validField(f reflect.StructField) bool {
	if isExported(f.Name) && f.Tag.Get("binary") != "ignore" {
		return true
	}
	return false
}

// isExported reports whether the identifier is exported.
func isExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}

//size of bool array when encode
func sizeofBoolArray(_len int) int {
	return SizeofUvarint(uint64(_len)) + (_len+8-1)/8
}

//size of string when encode
func sizeofString(_len int) int {
	return SizeofUvarint(uint64(_len)) + _len
}

//size of fix array, like []int16, []int64
func sizeofFixArray(_len, elemLen int) int {
	return SizeofUvarint(uint64(_len)) + _len*elemLen
}

func validUserType(t reflect.Type) bool {
	return sizeofNilPointer(t) >= 0
}
