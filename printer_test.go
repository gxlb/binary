package binary

import (
	"fmt"
	"reflect"
	"testing"
)

func _TestShowView(t *testing.T) {
	fmt.Printf("%#v\n", full.BaseStruct)
	fmt.Println(ShowString(full.BaseStruct))
	a := baseStruct{
		Int8:       18,
		Int16:      4660,
		Int32:      305419896,
		Int64:      1311768467463790320,
		Uint8:      18,
		Uint16:     4660,
		Uint32:     1898136936,
		Uint64:     11611200575286075120,
		Float32:    1234.567749,
		Float64:    2345.678901,
		Complex64:  complex(1.124565, 2.344565),
		Complex128: complex(333.456979, 567.345779),
		Array:      [4]uint8{1, 2, 3, 4},
		Bool:       false,
		BoolArray:  [9]bool{true, false, false, false, false, true, true, false, true},
	}
	fmt.Println(reflect.DeepEqual(a, full.BaseStruct))
	fmt.Printf("%#v\n", a)
	x := &[1]string{"x"}
	fmt.Printf("%#v\n", x)

	type A struct {
		m string
		n int
		o *[]string
	}
	type B struct {
		x []int
		y string
		z A
	}
	type T struct {
		a int
		b string
		c []string
		d map[string]B
		e B
	}
	var xx T
	fmt.Println(ShowString(xx))
}
