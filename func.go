package binary

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"
)

func assert(b bool, msg interface{}) {
	if !b {
		panic(fmt.Errorf("%s", msg))
	}
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
// It is recommended to use RegisterType to improve this case.
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

func validUserType(t reflect.Type) bool {
	return sizeofNilPointer(t) >= 0
}

func indirectType(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return t.Elem()
	}
	return t
}
