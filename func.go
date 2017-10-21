package binary

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

//var (
//	tSizer        reflect.Type
//	tPacker       reflect.Type
//	tUnpacker     reflect.Type
//	tPackUnpacker reflect.Type
//)

//func init() {
//	var sizer Sizer
//	var packer Packer
//	var unpacker Unpacker
//	var packUnpacker PackUnpacker
//	tSizer = reflect.TypeOf(&sizer).Elem()
//	tPacker = reflect.TypeOf(&packer).Elem()
//	tUnpacker = reflect.TypeOf(&unpacker).Elem()
//	tPackUnpacker = reflect.TypeOf(&packUnpacker).Elem()
//}

func sizeof(i interface{}) int {
	switch d := i.(type) { //fast size calculation
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
	}

	s := sizeofValue(reflect.ValueOf(i))
	return s
}

func assert(b bool, msg interface{}) {
	if !b {
		panic(msg)
	}
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeofValue(v reflect.Value) (l int) {
	if v.Kind() == reflect.Ptr && v.IsNil() { //nil is not aviable
		return sizeofNilPointer(v.Type())
	}

	v = reflect.Indirect(v)                 //redrect pointer to it's value
	if s := _fixTypeSize(v.Type()); s > 0 { //fixed size
		return s
	}
	switch t := v.Type(); t.Kind() {
	case reflect.Int:
		return SizeofVarint(v.Int())
	case reflect.Uint:
		return SizeofUvarint(v.Uint())
	case reflect.Slice, reflect.Array:
		arrayLen := v.Len()
		if s := _fixTypeSize(t.Elem()); s > 0 {
			if t.Elem().Kind() == reflect.Bool {
				return sizeofBoolArray(arrayLen)
			}
			return sizeofFixArray(arrayLen, s)
		} else {
			sum := SizeofUvarint(uint64(arrayLen)) //array size bytes num
			if sizeofNilPointer(t.Elem()) < 0 {    //check if array element type valid
				return -1
			}
			for i, n := 0, arrayLen; i < n; i++ {
				s := sizeofValue(v.Index(i))
				assert(s >= 0, v.Type().Kind().String()) //element size must not error
				sum += s
			}
			return sum
		}
	case reflect.Map:
		mapLen := v.Len()
		sum := SizeofUvarint(uint64(mapLen)) //array size
		keys := v.MapKeys()

		if sizeofNilPointer(t.Key()) < 0 ||
			sizeofNilPointer(t.Elem()) < 0 { //check if map key and value type valid
			return -1
		}

		for i := 0; i < mapLen; i++ {
			key := keys[i]
			sizeKey := sizeofValue(key)
			assert(sizeKey >= 0, key.Type().Kind().String()) //key size must not error

			sum += sizeKey
			value := v.MapIndex(key)
			sizeValue := sizeofValue(value)
			assert(sizeValue >= 0, value.Type().Kind().String()) //key size must not error

			sum += sizeValue
		}
		return sum

	case reflect.Struct:
		return queryStruct(v.Type()).sizeofValue(v)

	case reflect.String:
		return sizeofString(v.Len()) //string length and data
	}
	return -1
}

func sizeofNilPointer(t reflect.Type) int {
	tt := t
	if tt.Kind() == reflect.Ptr {
		tt = t.Elem()
	}
	if s := _fixTypeSize(tt); s > 0 { //fix size
		return s
	}
	switch tt.Kind() {
	case reflect.Int, reflect.Uint: //zero varint will be encoded as 1 byte
		return 1
	case reflect.Slice, reflect.String, reflect.Map:
		return SizeofUvarint(0)
	case reflect.Array:
		if s := _fixTypeSize(tt.Elem()); s > 0 {
			if tt.Elem().Kind() == reflect.Bool {
				return sizeofBoolArray(tt.Len())
			}
			return sizeofFixArray(tt.Len(), s)
		} else {
			size := sizeofNilPointer(tt.Elem())
			if size < 0 {
				return -1
			}
			return sizeofFixArray(tt.Len(), size)
		}

	case reflect.Struct:
		return queryStruct(tt).sizeofEmptyPointer(tt)
	}

	return -1
}

func _fixTypeSize(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
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
func newPtr(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr {
		e := v.Type().Elem()
		switch e.Kind() {
		case reflect.Int, reflect.Uint, reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16,
			reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64,
			reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64,
			reflect.Complex128, reflect.String, reflect.Array, reflect.Struct, reflect.Slice, reflect.Map:
			if v.IsNil() {
				v.Set(reflect.New(e))
			}
			return true
		}
	}
	return false
}

// NOTE:
// This function will make the encode/decode of struct slow down.
// It is recommeded to use RegStruct to improve this case.
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
