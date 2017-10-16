package binary

import (
	"reflect"
	"unicode"
	"unicode/utf8"
)

//var nameToType = map[string]reflect.Kind{
//	"uint8":  reflect.Uint8,
//	"uint16": reflect.Uint16,
//	"uint32": reflect.Uint32,
//	"uint64": reflect.Uint64,
//	"int8":   reflect.Int8,
//	"int16":  reflect.Int16,
//	"int32":  reflect.Int32,
//	"int64":  reflect.Int64,
//}

//func getKind(kind string) reflect.Kind {
//	if k, ok := nameToType[kind]; ok {
//		return k
//	}
//	return reflect.Invalid
//}

func sizeof(i interface{}) int {
	switch d := i.(type) { //fast size calculation
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case int16, uint16, *int16, *uint16:
		return 2
	case int32, uint32, *int32, *uint32, float32, *float32:
		return 4
	case int64, uint64, *int64, *uint64, float64, *float64, complex64, *complex64, int, *int, uint, *uint:
		return 8
	case complex128:
		return 16
	case string:
		return sizeofString(len(d))
	}

	s := sizeofValue(reflect.ValueOf(i))
	return s
}

// sizeof returns the size >= 0 of variables for the given type or -1 if the type is not acceptable.
func sizeofValue(v reflect.Value) (l int) {
	//	defer func() {
	//		fmt.Printf("sizeof(%s)=%d\n", v.Type().String(), l)
	//	}()
	if v.Kind() == reflect.Ptr && v.IsNil() { //nil is not aviable
		return -1
	}

	v = reflect.Indirect(v)                 //redrect pointer to it's value
	if s := _fixTypeSize(v.Type()); s > 0 { //fix size
		return s
	}
	switch t := v.Type(); t.Kind() {
	case reflect.Slice, reflect.Array:
		if s := _fixTypeSize(t.Elem()); s > 0 {
			if t.Elem().Kind() == reflect.Bool {
				return sizeofBoolArray(v.Len())
			}
			return SizeofUvarint(uint64(v.Len())) + s*v.Len()
		} else {
			sum := SizeofUvarint(uint64(v.Len())) //array size
			for i, n := 0, v.Len(); i < n; i++ {
				s := sizeofValue(v.Index(i))
				if s < 0 {
					return -1
				}
				sum += s
			}
			return sum
		}
	case reflect.Map:
		sum := SizeofUvarint(uint64(v.Len())) //array size
		keys := v.MapKeys()
		l := len(keys)
		for i := 0; i < l; i++ {
			key := keys[i]
			s := sizeofValue(key)
			if s < 0 {
				return -1
			}
			sum += s
			s2 := sizeofValue(v.MapIndex(key))
			if s < 0 {
				return -1
			}
			sum += s2
		}
		return sum

	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			if validField(v.Field(i), v.Type().Field(i)) {
				s := sizeofValue(v.Field(i))
				if s < 0 {
					return -1
				}
				sum += s
			}
		}
		return sum

	case reflect.String:
		return sizeofString(v.Len()) //string length and data
	}

	return -1
}

func sizeofEmptyValue(v reflect.Value) (l int) {
	if v.Kind() == reflect.Ptr && v.IsNil() { //nil is not aviable
		return -1
	}

	v = reflect.Indirect(v)                 //redrect fointer to it's value
	if s := _fixTypeSize(v.Type()); s > 0 { //fix size
		return s
	}
	switch t := v.Type(); t.Kind() {
	case reflect.Slice, reflect.Array, reflect.String:
		return SizeofUvarint(uint64(0))

	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeofEmptyValue(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
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
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.Complex64, reflect.Int, reflect.Uint:
		return 8
	case reflect.Complex128:
		return 16
	}
	return -1
}

func newPtr(v reflect.Value) bool {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		e := v.Type().Elem()
		switch e.Kind() {
		case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16,
			reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64,
			reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64,
			reflect.Complex128, reflect.String, reflect.Array, reflect.Struct, reflect.Slice, reflect.Map:
			v.Set(reflect.New(e))
		default:
			return false
		}
		return true
	}
	return false
}

func validField(v reflect.Value, f reflect.StructField) bool {
	//println("validField", v.CanSet(), f.Name, f.Index)
	if //v.CanSet() ||
	isExported(f.Name) && f.Tag.Get("binary") != "ignore" {
		return true
	}
	return false
}

// isExported reports whether the identifier is exported.
func isExported(id string) bool {
	r, _ := utf8.DecodeRuneInString(id)
	return unicode.IsUpper(r)
}

//deep indirect change ***X to X
//func DeepIndirect(v reflect.Value) reflect.Value {
//	for v.Kind() == reflect.Ptr {
//		v = v.Elem()
//	}
//	return v
//}

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
