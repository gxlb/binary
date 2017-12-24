package binary

import (
	"fmt"
	"io"
	"reflect"
	"testing"
	"time"
	"unsafe"
)

type TDoNotSupport struct {
	DeepPointer   **uint32
	Uintptr       uintptr
	UnsafePointer unsafe.Pointer
	Ch            chan bool
	Map           map[uintptr]uintptr
	Map2          map[int]uintptr
	Map3          map[uintptr]int
	Slice         []uintptr
	Array         [2]uintptr
	Array2        [2][2]uintptr
	Array3        [2]struct{ A uintptr }
	Func          func()
	Struct        struct {
		Uintptr uintptr
	}
	Struct2 struct {
		PStruct *struct {
			PPUintptr **uintptr
		}
	}
	Struct3 struct {
		PStruct *struct {
			PUintptr  *uintptr
			PPUintptr **uintptr
		}
	}
	PStruct *struct {
		PUintptr  *uintptr
		PPUintptr **uintptr
	}
}

var doNotSupportTypes = TDoNotSupport{
	Map2: map[int]uintptr{1: 1},
	Map3: map[uintptr]int{2: 2},
}

type fastValues struct {
	Int             int
	Uint            uint
	Bool            bool
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	String          string
	IntSlice        []int
	UintSlice       []uint
	BoolSlice       []bool
	Int8Slice       []int8
	Int16Slice      []int16
	Int32Slice      []int32
	Int64Slice      []int64
	Uint8Slice      []uint8
	Uint16Slice     []uint16
	Uint32Slice     []uint32
	Uint64Slice     []uint64
	Float32Slice    []float32
	Float64Slice    []float64
	Complex64Slice  []complex64
	Complex128Slice []complex128
	StringSlice     []string
}

var _fastValues = fastValues{
	Int:             -2,
	Uint:            2,
	Bool:            true,
	Int8:            -3,
	Int16:           -4,
	Int32:           -5,
	Int64:           -6,
	Uint8:           3,
	Uint16:          4,
	Uint32:          5,
	Uint64:          6,
	Float32:         -7,
	Float64:         7,
	Complex64:       8,
	Complex128:      9,
	String:          "hello",
	IntSlice:        []int{-1, 2},
	UintSlice:       []uint{1, 3},
	BoolSlice:       []bool{false, true},
	Int8Slice:       []int8{-1, 2},
	Int16Slice:      []int16{-1, 3},
	Int32Slice:      []int32{-1, 4},
	Int64Slice:      []int64{-1, 5},
	Uint8Slice:      []uint8{1, 6},
	Uint16Slice:     []uint16{1, 7},
	Uint32Slice:     []uint32{1, 8},
	Uint64Slice:     []uint64{1, 9},
	Float32Slice:    []float32{1, 10.1},
	Float64Slice:    []float64{1, 11.2},
	Complex64Slice:  []complex64{1, 2.2},
	Complex128Slice: []complex128{1, 12.9},
	StringSlice:     []string{"abc", "bcd"},
}

type baseStruct struct {
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	Array      [4]uint8
	Bool       bool
	BoolArray  [9]bool
}

type littleStruct struct {
	String string
	Int16  int16
}

type fullStruct struct {
	BaseStruct    baseStruct
	LittleStruct  littleStruct
	PLittleStruct *littleStruct
	String        string
	PString       *string
	PInt32        *int32
	Slice         []*littleStruct
	PSlice        *[]*string
	Float64Slice  []float64
	BoolSlice     []bool
	Uint32Slice   []uint32
	Map           map[string]*littleStruct
	Map2          map[string]uint16
	IntSlice      []int
	UintSlice     []uint
}

var full = fullStruct{
	BaseStruct: baseStruct{
		Int8:       0x12,
		Int16:      0x1234,
		Int32:      0x12345678,
		Int64:      0x123456789abcdef0,
		Uint8:      0x12,
		Uint16:     0x1234,
		Uint32:     0x71234568,
		Uint64:     0xa123456789bcdef0,
		Float32:    1234.5678,
		Float64:    2345.6789012,
		Complex64:  complex(1.12456453, 2.344565),
		Complex128: complex(333.4569789789123, 567.34577890012),
		Array:      [4]uint8{0x1, 0x2, 0x3, 0x4},
		Bool:       false,
		BoolArray:  [9]bool{true, false, false, false, false, true, true, false, true},
	},
	LittleStruct: littleStruct{
		String: "abc",
		Int16:  0x1234,
	},
	PLittleStruct: &littleStruct{
		String: "bcd",
		Int16:  0x2345,
	},
	String:  "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	PString: newString("hello"),
	PInt32:  newInt32(0x11223344),
	Slice: []*littleStruct{
		&littleStruct{
			String: "abc",
			Int16:  0x1122,
		},
		&littleStruct{
			String: "bcd",
			Int16:  0x2233,
		},
		&littleStruct{
			String: "cdef",
			Int16:  0x3344,
		},
	},
	PSlice:       &[]*string{newString("abc"), newString("def"), newString("ghijkl")},
	Float64Slice: []float64{3.141592654, 1.137856998, 6.789012345},
	BoolSlice:    []bool{false, true, false, false, true, true, false},
	Uint32Slice:  []uint32{0x12345678, 0x23456789, 0x34567890, 0x4567890a, 0x567890ab},
	Map:          map[string]*littleStruct{"a": &littleStruct{String: "a", Int16: 0x1122}, "b": &littleStruct{String: "b", Int16: 0x1122}},
	Map2:         map[string]uint16{"aaa": 0x5566, "bbb": 0x7788},
	IntSlice:     []int{0, -1, 1, -2, 2, -63, 63, -64, 64, -65, 65, -125, 125, -126, 126, -127, 127, -128, 128, -32765, 32765, -32766, 32766, -32767, 32767, -32768, 32768, -2147483645, 2147483645, -2147483646, 2147483646, -2147483647, 2147483647, -2147483648, 2147483648, -9223372036854775807, 9223372036854775806, -9223372036854775808, 9223372036854775807},
	UintSlice:    []uint{0, 1, 2, 127, 128, 32765, 32766, 32767, 32768, 65533, 65534, 65535, 65536, 0xFFFFFD, 0xFFFFFE, 0xFFFFFF, 0xFFFFFFFFFFFFFFFD, 0xFFFFFFFFFFFFFFFE, 0xFFFFFFFFFFFFFFFF},
}

func newString(s string) *string {
	p := new(string)
	*p = s
	return p
}
func newInt32(i int32) *int32 {
	p := new(int32)
	*p = i
	return p
}

var bigFull = []byte{
	//field#1|BaseStruct|binary.baseStruct{Int8:18, Int16:4660, Int32:305419896, Int64:1311768467463790320, Uint8:0x12, Uint16:0x1234, Uint32:0x71234568, Uint64:0xa123456789bcdef0, Float32:1234.5677, Float64:2345.6789012, Complex64:(1.1245645+2.344565i), Complex128:(333.4569789789123+567.34577890012i), Array:[4]uint8{0x1, 0x2, 0x3, 0x4}, Bool:false, BoolArray:[9]bool{true, false, false, false, false, true, true, false, true}}
	0x12, 0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x12, 0x34, 0x71, 0x23, 0x45, 0x68, 0xa1, 0x23, 0x45, 0x67, 0x89, 0xbc, 0xde, 0xf0, 0x44, 0x9a, 0x52, 0x2b, 0x40, 0xa2, 0x53, 0x5b, 0x98, 0xf0, 0x26, 0x6e, 0x3f, 0x8f, 0xf1, 0xbb, 0x40, 0x16, 0x0d, 0x5a, 0x40, 0x74, 0xd7, 0x4f, 0xc9, 0x30, 0x96, 0x34, 0x40, 0x81, 0xba, 0xc4, 0x27, 0xba, 0x5d, 0x4c, 0x04, 0x01, 0x02, 0x03, 0x04, 0x00, 0x09, 0x61, 0x01,
	//field#2|LittleStruct|binary.littleStruct{String:"abc", Int16:4660}
	0x03, 0x61, 0x62, 0x63, 0x12, 0x34,
	//field#3|PLittleStruct|&binary.littleStruct{String:"bcd", Int16:9029}
	0x03, 0x62, 0x63, 0x64, 0x23, 0x45,
	//field#4|String|"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
	//field#5|PString|(*string)(0xc042033170)
	0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	//field#6|PInt32|(*int32)(0xc04200e268)
	0x11, 0x22, 0x33, 0x44,
	//field#7|Slice|[]*binary.littleStruct{(*binary.littleStruct)(0x68f620), (*binary.littleStruct)(0x68f640), (*binary.littleStruct)(0x68f660)}
	0x03, 0x07, 0x03, 0x61, 0x62, 0x63, 0x11, 0x22, 0x03, 0x62, 0x63, 0x64, 0x22, 0x33, 0x04, 0x63, 0x64, 0x65, 0x66, 0x33, 0x44,
	//field#8|PSlice|&[]*string{(*string)(0xc042033180), (*string)(0xc042033190), (*string)(0xc0420331a0)}
	0x03, 0x07, 0x03, 0x61, 0x62, 0x63, 0x03, 0x64, 0x65, 0x66, 0x06, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c,
	//field#9|Float64Slice|[]float64{3.141592654, 1.137856998, 6.789012345}
	0x03, 0x40, 0x09, 0x21, 0xfb, 0x54, 0x52, 0x45, 0x50, 0x3f, 0xf2, 0x34, 0xa9, 0x8a, 0x1e, 0xf4, 0xaf, 0x40, 0x1b, 0x27, 0xf2, 0xda, 0x27, 0xa9, 0x3c,
	//field#10|BoolSlice|[]bool{false, true, false, false, true, true, false}
	0x07, 0x32,
	//field#11|Uint32Slice|[]uint32{0x12345678, 0x23456789, 0x34567890, 0x4567890a, 0x567890ab}
	0x05, 0x12, 0x34, 0x56, 0x78, 0x23, 0x45, 0x67, 0x89, 0x34, 0x56, 0x78, 0x90, 0x45, 0x67, 0x89, 0x0a, 0x56, 0x78, 0x90, 0xab,
	//field#12|Map|map[string]*binary.littleStruct{"a":(*binary.littleStruct)(0xc0420029c0), "b":(*binary.littleStruct)(0xc0420029e0)}
	0x02, 0x01, 0x61, 0x03, 0x01, 0x61, 0x11, 0x22, 0x01, 0x62, 0x01, 0x62, 0x11, 0x22,
	//field#13|Map2|map[string]uint16{"aaa":0x5566, "bbb":0x7788}
	0x02, 0x03, 0x61, 0x61, 0x61, 0x55, 0x66, 0x03, 0x62, 0x62, 0x62, 0x77, 0x88,
	//field#14|IntSlice|[]int{0, -1, 1, -2, 2, -63, 63, -64, 64, -65, 65, -125, 125, -126, 126, -127, 127, -128, 128, -32765, 32765, -32766, 32766, -32767, 32767, -32768, 32768, -2147483645, 2147483645, -2147483646, 2147483646, -2147483647, 2147483647, -2147483648, 2147483648, -9223372036854775807, 9223372036854775806, -9223372036854775808, 9223372036854775807}
	0x27, 0x00, 0x01, 0x02, 0x03, 0x04, 0x7d, 0x7e, 0x7f, 0x80, 0x01, 0x81, 0x01, 0x82, 0x01, 0xf9, 0x01, 0xfa, 0x01, 0xfb, 0x01, 0xfc, 0x01, 0xfd, 0x01, 0xfe, 0x01, 0xff, 0x01, 0x80, 0x02, 0xf9, 0xff, 0x03, 0xfa, 0xff, 0x03, 0xfb, 0xff, 0x03, 0xfc, 0xff, 0x03, 0xfd, 0xff, 0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03, 0x80, 0x80, 0x04, 0xf9, 0xff, 0xff, 0xff, 0x0f, 0xfa, 0xff, 0xff, 0xff, 0x0f, 0xfb, 0xff, 0xff, 0xff, 0x0f, 0xfc, 0xff, 0xff, 0xff, 0x0f, 0xfd, 0xff, 0xff, 0xff, 0x0f, 0xfe, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0x0f, 0x80, 0x80, 0x80, 0x80, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
	//field#15|UintSlice|[]uint{0x0, 0x1, 0x2, 0x7f, 0x80, 0x7ffd, 0x7ffe, 0x7fff, 0x8000, 0xfffd, 0xfffe, 0xffff, 0x10000, 0xfffffd, 0xfffffe, 0xffffff, 0xfffffffffffffffd, 0xfffffffffffffffe, 0xffffffffffffffff}
	0x13, 0x00, 0x01, 0x02, 0x7f, 0x80, 0x01, 0xfd, 0xff, 0x01, 0xfe, 0xff, 0x01, 0xff, 0xff, 0x01, 0x80, 0x80, 0x02, 0xfd, 0xff, 0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03, 0x80, 0x80, 0x04, 0xfd, 0xff, 0xff, 0x07, 0xfe, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0x07, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
}

var littleFull = []byte{
	//field#1|BaseStruct|binary.baseStruct{Int8:18, Int16:4660, Int32:305419896, Int64:1311768467463790320, Uint8:0x12, Uint16:0x1234, Uint32:0x71234568, Uint64:0xa123456789bcdef0, Float32:1234.5677, Float64:2345.6789012, Complex64:(1.1245645+2.344565i), Complex128:(333.4569789789123+567.34577890012i), Array:[4]uint8{0x1, 0x2, 0x3, 0x4}, Bool:false, BoolArray:[9]bool{true, false, false, false, false, true, true, false, true}}
	0x12, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12, 0x12, 0x34, 0x12, 0x68, 0x45, 0x23, 0x71, 0xf0, 0xde, 0xbc, 0x89, 0x67, 0x45, 0x23, 0xa1, 0x2b, 0x52, 0x9a, 0x44, 0x6e, 0x26, 0xf0, 0x98, 0x5b, 0x53, 0xa2, 0x40, 0xbb, 0xf1, 0x8f, 0x3f, 0x5a, 0x0d, 0x16, 0x40, 0x34, 0x96, 0x30, 0xc9, 0x4f, 0xd7, 0x74, 0x40, 0x4c, 0x5d, 0xba, 0x27, 0xc4, 0xba, 0x81, 0x40, 0x04, 0x01, 0x02, 0x03, 0x04, 0x00, 0x09, 0x61, 0x01,
	//field#2|LittleStruct|binary.littleStruct{String:"abc", Int16:4660}
	0x03, 0x61, 0x62, 0x63, 0x34, 0x12,
	//field#3|PLittleStruct|&binary.littleStruct{String:"bcd", Int16:9029}
	0x03, 0x62, 0x63, 0x64, 0x45, 0x23,
	//field#4|String|"0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	0x40, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
	//field#5|PString|(*string)(0xc042033170)
	0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f,
	//field#6|PInt32|(*int32)(0xc04200e268)
	0x44, 0x33, 0x22, 0x11,
	//field#7|Slice|[]*binary.littleStruct{(*binary.littleStruct)(0x68f400), (*binary.littleStruct)(0x68f420), (*binary.littleStruct)(0x68f440)}
	0x03, 0x07, 0x03, 0x61, 0x62, 0x63, 0x22, 0x11, 0x03, 0x62, 0x63, 0x64, 0x33, 0x22, 0x04, 0x63, 0x64, 0x65, 0x66, 0x44, 0x33,
	//field#8|PSlice|&[]*string{(*string)(0xc042033180), (*string)(0xc042033190), (*string)(0xc0420331a0)}
	0x03, 0x07, 0x03, 0x61, 0x62, 0x63, 0x03, 0x64, 0x65, 0x66, 0x06, 0x67, 0x68, 0x69, 0x6a, 0x6b, 0x6c,
	//field#9|Float64Slice|[]float64{3.141592654, 1.137856998, 6.789012345}
	0x03, 0x50, 0x45, 0x52, 0x54, 0xfb, 0x21, 0x09, 0x40, 0xaf, 0xf4, 0x1e, 0x8a, 0xa9, 0x34, 0xf2, 0x3f, 0x3c, 0xa9, 0x27, 0xda, 0xf2, 0x27, 0x1b, 0x40,
	//field#10|BoolSlice|[]bool{false, true, false, false, true, true, false}
	0x07, 0x32,
	//field#11|Uint32Slice|[]uint32{0x12345678, 0x23456789, 0x34567890, 0x4567890a, 0x567890ab}
	0x05, 0x78, 0x56, 0x34, 0x12, 0x89, 0x67, 0x45, 0x23, 0x90, 0x78, 0x56, 0x34, 0x0a, 0x89, 0x67, 0x45, 0xab, 0x90, 0x78, 0x56,
	//field#12|Map|map[string]*binary.littleStruct{"b":(*binary.littleStruct)(0xc0420029e0), "a":(*binary.littleStruct)(0xc0420029c0)}
	0x02, 0x01, 0x61, 0x03, 0x01, 0x61, 0x22, 0x11, 0x01, 0x62, 0x01, 0x62, 0x22, 0x11,
	//field#13|Map2|map[string]uint16{"aaa":0x5566, "bbb":0x7788}
	0x02, 0x03, 0x61, 0x61, 0x61, 0x66, 0x55, 0x03, 0x62, 0x62, 0x62, 0x88, 0x77,
	//field#14|IntSlice|[]int{0, -1, 1, -2, 2, -63, 63, -64, 64, -65, 65, -125, 125, -126, 126, -127, 127, -128, 128, -32765, 32765, -32766, 32766, -32767, 32767, -32768, 32768, -2147483645, 2147483645, -2147483646, 2147483646, -2147483647, 2147483647, -2147483648, 2147483648, -9223372036854775807, 9223372036854775806, -9223372036854775808, 9223372036854775807}
	0x27, 0x00, 0x01, 0x02, 0x03, 0x04, 0x7d, 0x7e, 0x7f, 0x80, 0x01, 0x81, 0x01, 0x82, 0x01, 0xf9, 0x01, 0xfa, 0x01, 0xfb, 0x01, 0xfc, 0x01, 0xfd, 0x01, 0xfe, 0x01, 0xff, 0x01, 0x80, 0x02, 0xf9, 0xff, 0x03, 0xfa, 0xff, 0x03, 0xfb, 0xff, 0x03, 0xfc, 0xff, 0x03, 0xfd, 0xff, 0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03, 0x80, 0x80, 0x04, 0xf9, 0xff, 0xff, 0xff, 0x0f, 0xfa, 0xff, 0xff, 0xff, 0x0f, 0xfb, 0xff, 0xff, 0xff, 0x0f, 0xfc, 0xff, 0xff, 0xff, 0x0f, 0xfd, 0xff, 0xff, 0xff, 0x0f, 0xfe, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0x0f, 0x80, 0x80, 0x80, 0x80, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
	//field#15|UintSlice|[]uint{0x0, 0x1, 0x2, 0x7f, 0x80, 0x7ffd, 0x7ffe, 0x7fff, 0x8000, 0xfffd, 0xfffe, 0xffff, 0x10000, 0xfffffd, 0xfffffe, 0xffffff, 0xfffffffffffffffd, 0xfffffffffffffffe, 0xffffffffffffffff}
	0x13, 0x00, 0x01, 0x02, 0x7f, 0x80, 0x01, 0xfd, 0xff, 0x01, 0xfe, 0xff, 0x01, 0xff, 0xff, 0x01, 0x80, 0x80, 0x02, 0xfd, 0xff, 0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03, 0x80, 0x80, 0x04, 0xfd, 0xff, 0xff, 0x07, 0xfe, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0x07, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
}

var littleFullAll = []byte{
	0x12, 0x34, 0x12, 0x78, 0x56, 0x34, 0x12, 0xf0, 0xde, 0xbc, 0x9a, 0x78, 0x56,
	0x34, 0x12, 0x12, 0x34, 0x12, 0x68, 0x45, 0x23, 0x71, 0xf0, 0xde, 0xbc, 0x89,
	0x67, 0x45, 0x23, 0xa1, 0x2b, 0x52, 0x9a, 0x44, 0x6e, 0x26, 0xf0, 0x98, 0x5b,
	0x53, 0xa2, 0x40, 0xbb, 0xf1, 0x8f, 0x3f, 0x5a, 0x0d, 0x16, 0x40, 0x34, 0x96,
	0x30, 0xc9, 0x4f, 0xd7, 0x74, 0x40, 0x4c, 0x5d, 0xba, 0x27, 0xc4, 0xba, 0x81,
	0x40, 0x04, 0x01, 0x02, 0x03, 0x04, 0xfe, 0x09, 0x61, 0x01, 0x03, 0x61, 0x62,
	0x63, 0x34, 0x12, 0x03, 0x62, 0x63, 0x64, 0x45, 0x23, 0x40, 0x30, 0x31, 0x32,
	0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66,
	0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63,
	0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39,
	0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36,
	0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x05, 0x68, 0x65, 0x6c,
	0x6c, 0x6f, 0x44, 0x33, 0x22, 0x11, 0x03, 0x03, 0x61, 0x62, 0x63, 0x22, 0x11,
	0x03, 0x62, 0x63, 0x64, 0x33, 0x22, 0x04, 0x63, 0x64, 0x65, 0x66, 0x44, 0x33,
	0x03, 0x1f, 0x03, 0x61, 0x62, 0x63, 0x03, 0x64, 0x65, 0x66, 0x06, 0x67, 0x68,
	0x69, 0x6a, 0x6b, 0x6c, 0x03, 0x50, 0x45, 0x52, 0x54, 0xfb, 0x21, 0x09, 0x40,
	0xaf, 0xf4, 0x1e, 0x8a, 0xa9, 0x34, 0xf2, 0x3f, 0x3c, 0xa9, 0x27, 0xda, 0xf2,
	0x27, 0x1b, 0x40, 0x07, 0x32, 0x05, 0x78, 0x56, 0x34, 0x12, 0x89, 0x67, 0x45,
	0x23, 0x90, 0x78, 0x56, 0x34, 0x0a, 0x89, 0x67, 0x45, 0xab, 0x90, 0x78, 0x56,
	0x02, 0x01, 0x61, 0x01, 0x61, 0x22, 0x11, 0x01, 0x62, 0x01, 0x62, 0x22, 0x11,
	0x02, 0x03, 0x61, 0x61, 0x61, 0x66, 0x55, 0x03, 0x62, 0x62, 0x62, 0x88, 0x77,
	0x27, 0x00, 0x01, 0x02, 0x03, 0x04, 0x7d, 0x7e, 0x7f, 0x80, 0x01, 0x81, 0x01,
	0x82, 0x01, 0xf9, 0x01, 0xfa, 0x01, 0xfb, 0x01, 0xfc, 0x01, 0xfd, 0x01, 0xfe,
	0x01, 0xff, 0x01, 0x80, 0x02, 0xf9, 0xff, 0x03, 0xfa, 0xff, 0x03, 0xfb, 0xff,
	0x03, 0xfc, 0xff, 0x03, 0xfd, 0xff, 0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03,
	0x80, 0x80, 0x04, 0xf9, 0xff, 0xff, 0xff, 0x0f, 0xfa, 0xff, 0xff, 0xff, 0x0f,
	0xfb, 0xff, 0xff, 0xff, 0x0f, 0xfc, 0xff, 0xff, 0xff, 0x0f, 0xfd, 0xff, 0xff,
	0xff, 0x0f, 0xfe, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0x0f, 0x80,
	0x80, 0x80, 0x80, 0x10, 0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0x01, 0xfc, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0x01, 0x13, 0x00, 0x01, 0x02, 0x7f, 0x80, 0x01, 0xfd,
	0xff, 0x01, 0xfe, 0xff, 0x01, 0xff, 0xff, 0x01, 0x80, 0x80, 0x02, 0xfd, 0xff,
	0x03, 0xfe, 0xff, 0x03, 0xff, 0xff, 0x03, 0x80, 0x80, 0x04, 0xfd, 0xff, 0xff,
	0x07, 0xfe, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0x07, 0xfd, 0xff, 0xff, 0xff,
	0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01,
}

var (
	typeMap = make(map[reflect.Type]string)
	strMap  = make(map[string]string)
)

func TestRegType(t *testing.T) {
	tp := reflect.TypeOf(full)
	typeMap[tp] = tp.String()
	strMap[tp.String()] = tp.String()

	types := []reflect.Type{
		tp,
		reflect.TypeOf("hello"),
		reflect.TypeOf(full.BaseStruct),
		reflect.TypeOf(full.BoolSlice),
		reflect.TypeOf(full.Float64Slice),
		reflect.TypeOf(full.Map2),
		reflect.TypeOf(full.LittleStruct),
	}

	n := 100
	start := time.Now()
	find := 0
	for i := 0; i < n; i++ {
		for _, v := range types {
			if _, ok := typeMap[v]; ok {
				find++
			}
		}
	}
	dur := time.Now().Sub(start)
	fmt.Printf("find typeMap n=%d*%d %d cost=%s find=%d\n", n, len(types), n*len(types), dur.String(), find)

	//	start = time.Now()
	//	find = 0
	//	for i := 0; i < n; i++ {
	//		for _, v := range types {
	//			if _, ok := strMap[v.String()]; ok {
	//				find++
	//			}
	//		}
	//	}
	//	dur = time.Now().Sub(start)
	//	fmt.Printf("find strMap n=%d*%d %d cost=%s find=%d\n", n, len(types), n*len(types), dur.String(), find)
}

func TestEncode(t *testing.T) {
	v := reflect.ValueOf(full)
	vt := v.Type()
	n := v.NumField()
	check := littleFull
	for i := 0; i < n; i++ {
		if !validField(vt.Field(i)) {
			continue
		}
		b, err := Encode(v.Field(i).Interface(), nil)
		c := check[:len(b)]
		check = check[len(b):]
		if err != nil {
			t.Error(err)
		}
		//fmt.Printf("//field#%d|%s|%#v\n$% x,\n", i+1, vt.Field(i).Name, v.Field(i).Interface(), b)
		if vt.Field(i).Type.Kind() != reflect.Map && //map keys will be got as unspecified order, byte order may change but it doesn't matter
			!reflect.DeepEqual(b, c) {
			//fmt.Printf("%d %s\ngot%#v\n%need#v\n", i, vt.Field(i).Type.String(), b, c)
			t.Errorf("field %d %s got %+v\nneed %+v\n", i, vt.Field(i).Name, b, c)
		}
	}

	////map fields will case uncertain bytes order but it does't matter
	//b2, err := Encode(full, nil)
	//if err != nil {
	//	t.Error(err)
	//}
	//fmt.Printf("//% x,\n", b2)
	//if !reflect.DeepEqual(b2, littleFull) {
	//	t.Errorf("got %+v\nneed %+v\n", b2, littleFull)
	//}
}

func TestEncodeBig(t *testing.T) {
	v := reflect.ValueOf(full)
	vt := v.Type()
	n := v.NumField()
	check := bigFull
	for i := 0; i < n; i++ {
		if !validField(vt.Field(i)) {
			continue
		}
		size := Sizeof(v.Field(i).Interface())
		encoder := NewEncoderEndian(size, BigEndian)
		err := encoder.Value(v.Field(i).Interface())
		b := encoder.Buffer()
		c := check[:len(b)]
		check = check[len(b):]
		if err != nil {
			t.Error(err)
		}
		//fmt.Printf("//field#%d|%s|%#v\n$% x,\n", i+1, vt.Field(i).Name, v.Field(i).Interface(), b)
		if vt.Field(i).Type.Kind() != reflect.Map && //map keys will be got as unspecified order, byte order may change but it doesn't matter
			!reflect.DeepEqual(b, c) {
			//fmt.Printf("%d %s\ngot%#v\n%need#v\n", i, vt.Field(i).Type.String(), b, c)
			t.Errorf("field %d %s got %+v\nneed %+v\n", i, vt.Field(i).Name, b, c)
		}
	}
}

func TestDecode(t *testing.T) {
	var v fullStruct
	err := Decode(littleFullAll, &v)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(v, full) {
		t.Errorf("got %#v\nneed %#v\n", v, full)
	}
}

func TestReset(t *testing.T) {
	encoder := NewEncoder(100)
	encoder.Uint64(0x1122334455667788, false)
	encoder.String("0123456789abcdef")
	oldCheck := []byte{0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x10, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66}
	old := encoder.Buffer()
	l := encoder.Len()
	if len(old) != l {
		t.Errorf("encode len error: got %#v\nneed %#v\n", old, l)
	}

	if !reflect.DeepEqual(old, oldCheck) {
		t.Errorf("got %#v\nneed %#v\n", old, oldCheck)
	}
	encoder.Reset()
	var s struct {
		PString  *string
		PSlice   *[]int
		PArray   *[2]bool
		PArray2  *[2]struct{ X *string }
		PInt     *int32
		PStruct  *struct{ A int }
		PStruct2 *struct{ B *[]string }
	}
	err := encoder.Value(&s)
	if err != nil {
		t.Error(err)
	}
	_new := encoder.Buffer()
	l2 := encoder.Len()
	newCheck := []byte{0x0}
	if len(_new) != l2 {
		t.Errorf("encode len error: got %#v\nneed %#v\n", _new, l2)
	}
	if !reflect.DeepEqual(_new, newCheck) {
		t.Errorf("got %#v\nneed %#v\n", _new, newCheck)
	}
	if s := encoder.Skip(1); s < 0 {
		t.Errorf("got %#v\nneed %#v\n", s, 1)
	}
	if s := encoder.Skip(encoder.Cap()); s >= 0 {
		t.Errorf("got %#v\nneed %#v\n", s, -1)
	}
	r := encoder.reserve(0)
	if r != nil {
		t.Errorf("got %#v\nneed %#v\n", r, nil)
	}

	defer func() {
		if e := recover(); e == nil {
			t.Error("need panic but not")
		}
	}()

	if !encoder.ResizeBuffer(101) {
		t.Errorf("Decoder: have %v, want %v", false, true)
	}

	large := [100]complex128{}
	err2 := encoder.Value(&large)
	if err2 == nil {
		t.Errorf("got err=nil, need err=none-nil")
	} else {
		//println("info******", err2.Error())
	}

	r2 := encoder.reserve(100)
	if r2 != nil {
		t.Errorf("got %#v\nneed %#v\n", r2, nil)
	}
}

func TestEncodeEmptyPointer(t *testing.T) {
	var s struct {
		PString  *string
		PSlice   *[]int
		PArray   *[2]bool
		PArray2  *[2]struct{ X *string }
		PInt     *int32
		PStruct  *struct{ A int }
		PStruct2 *struct{ B *[]string }
	}
	b, err := Encode(&s, nil)
	if err != nil {
		t.Error(err)
	}
	ss := s

	err = Decode(b, &ss)
	if err != nil {
		t.Error(err)
	}

	b2, err2 := Encode(&ss, nil)
	if err2 != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, b2) {
		t.Errorf("%+v->%+v got %+v\nneed %+v\n", s, ss, b2, b)
	}
	check := []byte{0x0}
	if !reflect.DeepEqual(b2, check) {
		t.Errorf("got %+v\nneed %+v\n", b2, check)
	}
}

func TestHideStructField(t *testing.T) {
	type T struct {
		A uint32
		b uint32
		_ uint32
		C uint32 `binary:"ignore"`
	}
	var s T
	s.A = 0x11223344
	s.b = 0x22334455
	s.C = 0x33445566
	check := []byte{0x44, 0x33, 0x22, 0x11}
	b, err := Encode(s, nil)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(b, check) {
		t.Errorf("%T: got %x; want %x", s, b, check)
	}
	var ss, ssCheck T
	ssCheck.A = s.A
	err = Decode(b, &ss)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ss, ssCheck) {
		t.Errorf("%T: got %q; want %q", s, ss, ssCheck)
	}
}

func TestEndian(t *testing.T) {
	if LittleEndian.String() != "LittleEndian" ||
		LittleEndian.GoString() != "binary.LittleEndian" {
		t.Error("LittleEndian")
	}
	if BigEndian.String() != "BigEndian" ||
		BigEndian.GoString() != "binary.BigEndian" {
		t.Error("BigEndian")
	}
}

func TestByteReaderWriter(t *testing.T) {
	buff := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	reader := BytesReader(buff[:])
	r := make([]byte, len(buff)-1)
	n, err := reader.Read(r)
	if err != nil {
		t.Error(err)
	}
	if n != len(r) {
		t.Errorf("got %+v\nneed %+v\n", n, len(r))
	}
	if check := buff[:len(r)]; !reflect.DeepEqual(r, check) {
		t.Errorf("got %+v\nneed %+v\n", r, check)
	}
	n2, err2 := reader.Read(r)
	if n2 != 1 || err2 == nil {
		t.Errorf("got %d %+v\nneed %d %+v\n", n2, err2, 1, io.EOF)
	}

	wBuff := make([]byte, len(buff)+1)
	writer := BytesWriter(wBuff)
	w := buff[:]
	n3, err3 := writer.Write(w)
	if err3 != nil {
		t.Error(err3)
	}
	if n3 != len(w) {
		t.Errorf("got %+v\nneed %+v\n", n3, len(w))
	}
	n4, err4 := writer.Write(w)
	if n4 != 1 || err4 == nil {
		t.Errorf("got %d %+v\nneed %d %+v\n", n4, err4, 1, ErrNotEnoughSpace)
	}
	wCheck := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	if c := wBuff; !reflect.DeepEqual(c, wCheck) {
		t.Errorf("got %+v\nneed %+v\n", c, wCheck)
	}
}

func TestDecoderSkip(t *testing.T) {
	type skipedStruct struct {
		S         string
		I         int
		U         uint
		Map       map[uint32]uint32
		BoolArray [5]bool
		U16Array  [5]uint16
		StrArray  [5]string
		Struct    struct{ A uint8 }
		Bool1     bool
		Pointer   *uint32
		Bool2     bool
		Packed1   uint32 `binary:"packed"`
		Packed2   int64  `binary:"packed"`
	}
	RegsterType((*skipedStruct)(nil))

	var w [5]skipedStruct
	for i := len(w) - 1; i >= 0; i-- {
		w[i].S = fmt.Sprintf("%d", i)
		w[i].I = i
		w[i].U = uint(i)
		w[i].Map = map[uint32]uint32{uint32(i): uint32(i), uint32(i + 1): uint32(i + 1)}
		w[i].Struct.A = uint8(i)
		w[i].U16Array[i] = uint16(i)
		w[i].BoolArray[i] = true
		w[i].Bool2 = true
		w[i].Packed1 = uint32(i)
		w[i].Packed2 = int64(i * 2)
		if i%2 == 0 {
			w[i].Pointer = new(uint32)
		}
	}

	var r [3]skipedStruct
	b, err := Encode(&w, nil)
	if err != nil {
		t.Error(err)
	}

	err2 := Decode(b, &r)
	if err2 != nil {
		t.Error(err2)
	}
	for i := len(r) - 1; i >= 0; i-- {
		if !reflect.DeepEqual(w[i], r[i]) {
			t.Errorf("%d got %+v\nneed %+v\n", i, w[i], r[i])
		}
	}
}

//func TestDecoderSkipError(t *testing.T) {
//	//	defer func() {
//	//		if msg := recover(); msg == nil {
//	//			t.Fatal("did not panic")
//	//		} else {
//	//			fmt.Println(msg)
//	//		}
//	//	}()

//	bytes := []byte{2, 0, 0, 0, 0}
//	var dataDecode [0]uintptr
//	decoder := NewDecoder(bytes)
//	n := decoder.skipByType(reflect.TypeOf(dataDecode), false)
//	if n >= 0 {
//		t.Errorf("DecoderSkipError: have %d, want %d", n, -1)
//	} else {
//		//println(n)
//	}
//}

func TestFastValue(t *testing.T) {
	s := _fastValues
	v := reflect.ValueOf(s)
	encoder := NewEncoder(Size(s))
	for i := v.NumField() - 1; i >= 0; i-- {
		f := v.Field(i)
		if err := encoder.Value(f.Interface()); err != nil {
			t.Error(err)
		}
	}
	buffer := encoder.Buffer()

	var r fastValues
	vr := reflect.ValueOf(&r)
	vr = reflect.Indirect(vr)
	decoder := NewDecoder(buffer)
	for i := vr.NumField() - 1; i >= 0; i-- {
		f := vr.Field(i).Addr()
		oldSize := decoder.Len()
		if err := decoder.Value(f.Interface()); err != nil {
			t.Error(err)
		}
		size := decoder.Len() - oldSize
		assert(size == Sizeof(f.Interface()), "")
	}
	if !reflect.DeepEqual(r, s) {
		t.Errorf("got %+v\nneed %+v\n", r, s)
	}
}

func TestEncodeDonotSupportedType(t *testing.T) {
	ts := doNotSupportTypes
	if _, err := Encode(ts, nil); err == nil {
		t.Errorf("EncodeDonotSupportedType: have err == nil, want non-nil")
	}

	buff := make([]byte, 0)
	ecoder := NewEncoder(100)
	decoder := NewDecoder(buff)

	tv := reflect.Indirect(reflect.ValueOf(&ts))
	for i, n := 0, tv.NumField(); i < n; i++ {
		if _, err := Encode(tv.Field(i).Interface(), nil); err == nil {
			t.Errorf("EncodeDonotSupportedType.%v: have err == nil, want non-nil", tv.Field(i).Type())
		} else {
			//fmt.Println(err)
		}

		if err := ecoder.Value(tv.Field(i).Interface()); err == nil {
			t.Errorf("EncodeDonotSupportedType.%v: have err == nil, want non-nil", tv.Field(i).Type())
		} else {
			//fmt.Println(err)
		}

		if err := Decode(buff, tv.Field(i).Addr().Interface()); err == nil {
			t.Errorf("Decode DonotSupportedType.%v: have err == nil, want non-nil", tv.Field(i).Type())
		} else {
			//fmt.Printf("Decode error: %#v\n%s\n", tv.Field(i).Addr().Type().String(), err.Error())
		}

		if err := decoder.value(tv.Field(i), true, false); err == nil {
			t.Errorf("EncodeDonotSupportedType.%v: have err == nil, want non-nil", tv.Field(i).Type())
		} else {
			//fmt.Println(err)
		}
	}

	if queryStruct(tv.Type()).decode(decoder, tv) == nil {
		t.Errorf("decode DonotSupportedType.%v: have err == nil, want non-nil", tv.Type())
	}
}

func TestDecoder(t *testing.T) {
	buffer := []byte{}
	decoder := NewDecoder(buffer)
	got := decoder.Skip(0)
	if got != -1 {
		t.Errorf("Decoder: have %d, want %d", got, -1)
	}
	n := decoder.skipByType(reflect.TypeOf(uintptr(0)), false)
	if n != -1 {
		t.Errorf("Decoder: have %d, want %d", n, -1)
	}
}

func TestAssert(t *testing.T) {
	defer func() {
		if msg := recover(); msg == nil {
			t.Fatal("did not panic")
		}
	}()

	message := "it will panic"
	assert(false, message)
}

func TestRegStruct(t *testing.T) {
	type StructForReg struct {
		A int
		B uint `binary:"ignore"`
		C int  `binary:"int32"`
		d string
		_ int32
		F float32
		S struct {
			A int
			B string
		}
		S2 struct {
			A int
			B string
		}
		PS *struct {
			A int32
			B string
		}
	}
	RegsterType((*StructForReg)(nil))
	if err := RegsterType((*StructForReg)(nil)); err == nil { //duplicate regist
		t.Errorf("RegStruct: have err == nil, want non-nil")
	}
	var a = StructForReg{
		A: -5,
		B: 6,
		C: 7,
		d: "hello",
		F: 3.14,
	}
	a.S.A = 9
	a.S.B = "abc"
	b, err := Encode(&a, nil)
	if err != nil {
		t.Error(err)
	}

	var r StructForReg
	//fmt.Printf("%#v\n%#v\n", a, b)
	err = Decode(b, &r)
	if err != nil {
		t.Error(err)
	}
	c := a
	c.B = 0
	c.d = ""
	//r.PS = nil //BUG: how to encode nil pointer?
	if !reflect.DeepEqual(r, c) {
		t.Errorf("RegStruct got %+v\nneed %+v\n", r, c)
	}
}

func TestRegistStructUnsupported(t *testing.T) {
	err := RegsterType(int(0))
	if err == nil {
		t.Errorf("RegistStructUnsupported: have err == nil, want non-nil")
	}

	info := queryStruct(reflect.TypeOf(doNotSupportTypes))
	if info != nil {
		t.Errorf("RegistStructUnsupported: have info == %v, want nil", info)
	}
	numField := info.numField()
	if numField != 0 {
		t.Errorf("RegistStructUnsupported: have numField == %v, want 0", numField)
	}
	field := info.field(0)
	if field != nil {
		t.Errorf("RegistStructUnsupported: have info == %v, want nil", field)
	}
}

func TestDecodeUvarintOverflow(t *testing.T) {
	data := [][]byte{
		[]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x2},
		[]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x1, 0, 0},
	}
	var d uint
	for _, v := range data {
		decoder := NewDecoder(v)
		err := decoder.Value(&d)
		if err == nil {
			t.Errorf("DecodeUvarintOverflow: have err == nil, want none-nil")
		} else {
			//println(err.Error())
		}
	}
}

type sizerOnly struct{ A uint8 }

func (obj sizerOnly) Size() int { return 1 }

type encoderOnly struct{ B uint8 }

func (obj encoderOnly) Encode(buffer []byte) ([]byte, error) { return nil, nil }

type decoderOnly struct {
	C uint8
}

func (obj *decoderOnly) Decode(buffer []byte) error { return nil }

type sizeencoderOnly struct {
	sizerOnly
	encoderOnly
}
type sizedecoderOnly struct {
	sizerOnly
	decoderOnly
}
type encodedecoderOnly struct {
	encoderOnly
	decoderOnly
}
type fullSerializer struct {
	sizerOnly
	encoderOnly
	decoderOnly
}
type fullSerializerError struct {
	fullSerializer
}

func (obj *fullSerializerError) Decode(buffer []byte) error {
	return fmt.Errorf("expected error")
}

func TestBinarySerializer(t *testing.T) {
	var a sizerOnly
	var b encoderOnly
	var c decoderOnly
	var d sizeencoderOnly
	var e sizedecoderOnly
	var f encodedecoderOnly
	var g fullSerializerError
	var h fullSerializer

	testCase := func(data interface{}, testcase int) (info interface{}) {
		defer func() {

			_info := recover()
			if _info != nil && info == nil {
				info = _info
				//fmt.Println(info)
			}
		}()
		switch testcase {
		case 1:
			Sizeof(data)
		case 2:
			if _, err := Encode(data, nil); err != nil {
				info = err
			}

		case 3:
			buff := make([]byte, 1000)
			if err := Decode(buff, data); err != nil {
				info = err
			}
		case 4:
			encoder := NewEncoder(100)
			if err := encoder.Value(data); err != nil {
				info = err
			}
		}
		return
	}

	testCode := func(data interface{}) (info interface{}) {
		for i := 1; i <= 4; i++ {
			if _info := testCase(data, i); _info != nil && info == nil {
				info = _info
			}
		}
		return
	}

	if info := testCode(&a); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&b); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&c); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&d); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&e); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&f); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&g); info == nil {
		t.Errorf("BinarySerializer: have err == nil, want none-nil")
	}
	if info := testCode(&h); info != nil {
		t.Errorf("BinarySerializer: have err == %#v, want nil", info)
	}
}

func TestFastSizeof(t *testing.T) {
	type interSize struct {
		iter interface{}
		size int
	}
	var cases = []interSize{
		interSize{bool(false), 1},
		interSize{int8(0), 1},
		interSize{uint8(0), 1},
		interSize{int16(0), 2},
		interSize{uint16(0), 2},
		interSize{int32(0), 4},
		interSize{uint32(0), 4},
		interSize{int64(0), 8},
		interSize{uint64(0), 8},
		interSize{float32(0), 4},
		interSize{float64(0), 8},
		interSize{complex64(0), 8},
		interSize{complex128(0), 16},
		interSize{string("hello"), 6},

		interSize{int(0), 1},
		interSize{uint(0), 1},

		interSize{[]bool{false, false, true}, 2},
		interSize{[]int8{0}, 2},
		interSize{[]uint8{0}, 2},
		interSize{[]int16{0}, 3},
		interSize{[]uint16{0}, 3},
		interSize{[]int32{0}, 5},
		interSize{[]uint32{0}, 5},
		interSize{[]int64{0}, 9},
		interSize{[]uint64{0}, 9},
		interSize{[]float32{0}, 5},
		interSize{[]float64{0}, 9},
		interSize{[]complex64{0}, 9},
		interSize{[]complex128{0}, 17},
		interSize{[]string{"hello"}, 7},

		interSize{&[]bool{false, false, true}, 2},
		interSize{&[]int8{0}, 2},
		interSize{&[]uint8{0}, 2},
		interSize{&[]int16{0}, 3},
		interSize{&[]uint16{0}, 3},
		interSize{&[]int32{0}, 5},
		interSize{&[]uint32{0}, 5},
		interSize{&[]int64{0}, 9},
		interSize{&[]uint64{0}, 9},
		interSize{&[]float32{0}, 5},
		interSize{&[]float64{0}, 9},
		interSize{&[]complex64{0}, 9},
		interSize{&[]complex128{0}, 17},
		interSize{&[]string{"hello"}, 7},

		interSize{uintptr(0), -1},
		interSize{(*[]int)(nil), -1},
	}
	for i, v := range cases {
		s := fastSizeof(v.iter)
		if s != v.size {
			t.Errorf("%d %#v got %d need %d", i, v.iter, s, v.size)
		}
	}
}

func TestPackedInts(t *testing.T) {
	type packedInts struct {
		A int16    `binary:"packed"`
		B int32    `binary:"packed"`
		C int64    `binary:"packed"`
		D uint16   `binary:"packed"`
		E uint32   `binary:"packed"`
		F uint64   `binary:"packed"`
		G []uint64 `binary:"packed"`
	}
	var data = packedInts{1, 2, 3, 4, 5, 6, []uint64{7, 8, 9}}
	RegsterType((*packedInts)(nil))
	b, err := Encode(data, nil)
	if err != nil {
		t.Error(err)
	}
	if s := Sizeof(data); s != len(b) {
		t.Errorf("PackedInts got %+v %+v\nneed %+v\n", len(b), b, s)
	}
	check := []byte{0x2, 0x4, 0x6, 0x4, 0x5, 0x6, 0x3, 0x7, 0x8, 0x9}
	if !reflect.DeepEqual(b, check) {
		t.Errorf("PackedInts %#v\n got %+v\nneed %+v\n", data, b, check)
	}

	var dataDecode packedInts
	err = Decode(b, &dataDecode)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(dataDecode, data) {
		t.Errorf("PackedInts got %+v\nneed %+v\n", dataDecode, data)
	}
}

func TestBools(t *testing.T) {
	type boolset struct {
		A uint8   //0x11
		B bool    //true
		C uint8   //0x22
		D []bool  //[]bool{true, false, true}
		E bool    //true
		F *uint32 //false
		G bool    //true
		H uint8   //0x33
	}
	var data = boolset{
		0x11, true, 0x22, []bool{true, false, true}, true, nil, true, 0x33,
	}
	b, err := Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	size := Sizeof(data)
	if size != len(b) {
		fmt.Printf("EncodeBools got %#v %+v\nneed %+v\n", len(b), b, size)
	}
	check := []byte{0x11, 0xb, 0x22, 0x3, 0x5, 0x33}
	if !reflect.DeepEqual(b, check) {
		t.Errorf("EncodeBools %#v\n got %+v\nneed %+v\n", data, b, check)
	}

	var dataDecode boolset
	err = Decode(b, &dataDecode)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(dataDecode, data) {
		t.Errorf("EncodeBools got %+v\nneed %+v\n", dataDecode, data)
	}
}
