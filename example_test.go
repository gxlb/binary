// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"

	std "encoding/binary"

	"github.com/vipally/binary"
)

func ExampleWrite() {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("%#v", buf.Bytes())

	// Output:
	// []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x9, 0x40}
}

func ExampleWrite_multi() {
	buf := new(bytes.Buffer)
	var data = []interface{}{
		uint16(61374),
		int8(-54),
		uint8(254),
	}
	for _, v := range data {
		err := binary.Write(buf, binary.LittleEndian, v)
		if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
	}
	fmt.Printf("%#v", buf.Bytes())

	// Output:
	// []byte{0xbe, 0xef, 0xca, 0xfe}
}

func ExampleRead() {
	var pi float64
	b := []byte{0x18, 0x2d, 0x44, 0x54, 0xfb, 0x21, 0x09, 0x40}
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &pi)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Print(pi)

	// Output:
	// 3.141592653589793
}

func ExampleEncode() {
	var data struct {
		A uint32
		B int
		C string
	}
	data.A = 0x11223344
	data.B = -5
	data.C = "hello"

	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println("binary.Encode failed:", err)
	}
	fmt.Printf("Encode:\n%+v\n%#v", data, b)

	// Output:
	// Encode:
	// {A:287454020 B:-5 C:hello}
	// []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleEncode_withbuffer() {
	var data struct {
		A uint32
		B int
		C string
	}
	data.A = 0x11223344
	data.B = -5
	data.C = "hello"
	size := binary.Size(data)
	buffer := make([]byte, size)

	b, err := binary.Encode(data, buffer)
	if err != nil {
		fmt.Println("binary.Encode failed:", err)
	}
	fmt.Printf("Encode:\n%+v\n%#v", data, b)

	// Output:
	// Encode:
	// {A:287454020 B:-5 C:hello}
	// []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleEncode_bools() {
	type boolset struct {
		A uint8   //0xa
		B bool    //true
		C uint8   //0xc
		D []bool  //[]bool{true, false, true}
		E bool    //true
		F *uint32 //false
		G bool    //true
		H uint8   //0x8
	}
	var data = boolset{
		0xa, true, 0xc, []bool{true, false, true}, true, nil, true, 0x8,
	}
	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Size(data); size != len(b) {
		fmt.Printf("Encode got %#v %+v\nneed %+v\n", len(b), b, size)
	}
	fmt.Printf("Encode bools:\n%+v\nsize=%d result=%#v", data, len(b), b)

	// Output:
	// Encode bools:
	// {A:10 B:true C:12 D:[true false true] E:true F:<nil> G:true H:8}
	// size=6 result=[]byte{0xa, 0xb, 0xc, 0x3, 0x5, 0x8}
}

func ExampleEncode_boolArray() {
	var data = []bool{true, true, true, false, true, true, false, false, true}
	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Size(data); size != len(b) {
		fmt.Printf("Encode bool array:\ngot %#v %+v\nneed %+v\n", len(b), b, size)
	}
	fmt.Printf("Encode bool array:\n%#v\nsize=%d result=%#v", data, len(b), b)

	// Output:
	// Encode bool array:
	// []bool{true, true, true, false, true, true, false, false, true}
	// size=3 result=[]byte{0x9, 0x37, 0x1}
}

func ExampleEncode_packedInts() {
	type regedPackedInts struct {
		A int16    `binary:"packed"`
		B int32    `binary:"packed"`
		C int64    `binary:"packed"`
		D uint16   `binary:"packed"`
		E uint32   `binary:"packed"`
		F uint64   `binary:"packed"`
		G []uint64 `binary:"packed"`
		H uint     `binary:"ignore"`
	}
	binary.RegisterType((*regedPackedInts)(nil))

	var data = regedPackedInts{1, 2, 3, 4, 5, 6, []uint64{7, 8, 9}, 10}
	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Size(data); size != len(b) {
		fmt.Printf("PackedInts got %+v %+v\nneed %+v\n", len(b), b, size)
	}

	fmt.Printf("Encode packed ints:\n%+v\nsize=%d result=%#v", data, len(b), b)

	// Output:
	// Encode packed ints:
	// {A:1 B:2 C:3 D:4 E:5 F:6 G:[7 8 9] H:10}
	// size=10 result=[]byte{0x2, 0x4, 0x6, 0x4, 0x5, 0x6, 0x3, 0x7, 0x8, 0x9}
}

func ExampleDecode() {
	var s struct {
		A uint32
		B int
		C string
	}
	buffer := []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

	err := binary.Decode(buffer, &s)
	if err != nil {
		fmt.Println("binary.Decode failed:", err)
	}
	fmt.Printf("%+v", s)

	// Output:
	// {A:287454020 B:-5 C:hello}
}
func ExampleEncoder() {
	encoder := binary.NewEncoder(100)

	encoder.Uint32(0x11223344, false)
	encoder.Varint(-5)
	encoder.String("hello")
	encodeResult := encoder.Buffer()
	fmt.Printf("%#v", encodeResult)

	// Output:
	// []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleDecoder() {
	buffer := []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

	decoder := binary.NewDecoder(buffer)
	u32 := decoder.Uint32(false)
	i, _ := decoder.Varint()
	str := decoder.String()
	fmt.Printf("%#v %#v %#v", u32, i, str)

	// Output:
	// 0x11223344 -5 "hello"
}

type S struct {
	A uint32
	B int
	C string
}

func (this S) Size() int {
	size := binary.Size(this.A) + binary.Size(this.C) + binary.Size(int16(this.B))
	return size
}
func (this S) Encode(buffer []byte) ([]byte, error) {
	buff, err := binary.MakeEncodeBuffer(this, buffer)
	if err != nil {
		return nil, err
	}
	encoder := binary.NewEncoderBuffer(buff)
	encoder.Value(this.A)
	encoder.Int16(int16(this.B), false)
	encoder.Value(this.C)
	return encoder.Buffer(), nil
}
func (this *S) Decode(buffer []byte) error {
	decoder := binary.NewDecoder(buffer)
	decoder.Value(&this.A)
	this.B = int(decoder.Int16(false))
	decoder.Value(&this.C)
	return nil
}

func ExampleBinarySerializer() {
	/*
		type S struct {
			A uint32
			B int
			C string
		}
		func (this S) Size() int {
			size := binary.Sizeof(this.A) + binary.Sizeof(this.C) + binary.Sizeof(int16(this.B))
			return size
		}
		func (this S) Encode() ([]byte, error) {
			encoder := binary.NewEncoder(this.Size())
			encoder.Value(this.A)
			encoder.Int16(int16(this.B), false)
			encoder.Value(this.C)
			return encoder.Buffer(), nil
		}
		func (this *S) Decode(buffer []byte) error {
			decoder := binary.NewDecoder(buffer)
			decoder.Value(&this.A)
			this.B = int(decoder.Int16(false))
			decoder.Value(&this.C)
			return nil
		}
	*/
	var data, dataDecode S
	data.A = 0x11223344
	data.B = -5
	data.C = "hello"

	if err := binary.RegisterType((*S)(nil)); err != nil {
		fmt.Println(err)
	}
	b, err := binary.EncodeX(&data, nil, true)
	if err != nil {
		fmt.Println("binary.Encode failed:", err)
	}

	err = binary.DecodeX(b, &dataDecode, true)
	if err != nil {
		fmt.Println("binary.Decode failed:", err)
	}
	fmt.Printf("%+v\n%#v\n%+v", data, b, dataDecode)

	// Output:
	// {A:287454020 B:-5 C:hello}
	// []byte{0x44, 0x33, 0x22, 0x11, 0xfb, 0xff, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	// {A:287454020 B:-5 C:hello}
}

func ExampleSizeof() {
	var s struct {
		Int8        int8
		Int16       int16
		Int32       int32
		Int64       int64
		Uint8       uint8
		Uint16      uint16
		Uint32      uint32
		Uint64      uint64
		Float32     float32
		Float64     float64
		Complex64   complex64
		Complex128  complex128
		Array       [10]uint8
		Bool        bool
		BoolArray   [100]bool
		Uint32Array [10]uint32
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	coder := gob.NewEncoder(buf)
	coder.Encode(s)
	gobSize := len(buf.Bytes())

	stdSize := std.Size(s)
	size := binary.Size(s)

	fmt.Printf("Sizeof(s)  = %d\nstd Size(s)= %d\ngob Size(s)= %d", size, stdSize, gobSize)

	// Output:
	// Sizeof(s)  = 133
	// std Size(s)= 217
	// gob Size(s)= 412
}

func ExampleRegisterType() {
	type someRegedStruct struct {
		A int    `binary:"ignore"`
		B uint64 `binary:"packed"`
		C string
		D uint
	}
	binary.RegisterType((*someRegedStruct)(nil))

	var data = someRegedStruct{1, 2, "hello", 3}
	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Size(data); size != len(b) {
		fmt.Printf("RegedStruct got %+v %+v\nneed %+v\n", len(b), b, size)
	}

	fmt.Printf("Encode reged struct:\n%+v\nsize=%d result=%#v", data, len(b), b)

	// Output:
	// Encode reged struct:
	// {A:1 B:2 C:hello D:3}
	// size=8 result=[]byte{0x2, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x3}
}

//func ExampleShowString() {
//	type X struct {
//		A int
//		B string
//	}
//	type Y struct {
//		D X
//		E string
//	}
//	type Z struct {
//		G Y
//		H string
//		I []string
//		J map[string]int
//	}
//	var z = Z{
//		G: Y{
//			D: X{
//				A: 1,
//				B: `"a"=1`,
//			},
//			E: `"b"=2`,
//		},
//		H: `zzz`,
//		I: []string{
//			`c:\x\y\z`,
//			`d:\a\b\c`,
//		},
//		J: map[string]int{
//			`abc`: 1,
//		},
//	}
//	fmt.Println(binary.ShowString(z))
//	// Output:
//	// binary_test.Z{
//	//     G: binary_test.Y{
//	//         D: binary_test.X{
//	//             A: 1,
//	//             B: `"a"=1`,
//	//         },
//	//         E: `"b"=2`,
//	//     },
//	//     H: `zzz`,
//	//     I: []string{
//	//         `c:\x\y\z`,
//	//         `d:\a\b\c`,
//	//     },
//	//     J: map[string]int{
//	//         `abc`: 1,
//	//     },
//	// }
//}
//func ExampleShowSingleLineString() {
//	type X struct {
//		A int
//		B string
//	}
//	type Y struct {
//		D X
//		E string
//	}
//	type Z struct {
//		G Y
//		H string
//		I []string
//		J map[string]int
//	}
//	var z = Z{
//		G: Y{
//			D: X{
//				A: 1,
//				B: `"a"=1`,
//			},
//			E: `"b"=2`,
//		},
//		H: `zzz`,
//		I: []string{
//			`c:\x\y\z`,
//			`d:\a\b\c`,
//		},
//		J: map[string]int{
//			`abc`: 1,
//		},
//	}
//	fmt.Println(binary.ShowSingleLineString(z))

//	// Output:
//	// binary_test.Z{G: binary_test.Y{D: binary_test.X{A: 1,B: `"a"=1`,},E: `"b"=2`,},H: `zzz`,I: []string{`c:\x\y\z`,`d:\a\b\c`,},J: map[string]int{`abc`: 1,},}
//	//
//}
