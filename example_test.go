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
	fmt.Printf("%#v", b)

	// Output:
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
	size := binary.Sizeof(data)
	buffer := make([]byte, size)

	b, err := binary.Encode(data, buffer)
	if err != nil {
		fmt.Println("binary.Encode failed:", err)
	}
	fmt.Printf("%#v", b)

	// Output:
	// []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleEncode_bools() {
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
	b, err := binary.Encode(data, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Sizeof(data); size != len(b) {
		fmt.Printf("Encode got %#v %+v\nneed %+v\n", len(b), b, size)
	}
	fmt.Printf("Encode %#v\nsize=%d result=%#v", data, len(b), b)

	// Output:
	// Encode binary_test.boolset{A:0x11, B:true, C:0x22, D:[]bool{true, false, true}, E:true, F:(*uint32)(nil), G:true, H:0x33}
	// size=6 result=[]byte{0x11, 0xb, 0x22, 0x3, 0x5, 0x33}
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

func (this *S) Size() int {
	size := binary.Sizeof(this.A) + binary.Sizeof(this.C) + binary.Sizeof(int16(this.B))
	return size
}
func (this *S) Encode(buffer []byte) ([]byte, error) {
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
		func (this *S) Size() int {
			size := binary.Sizeof(this.A) + binary.Sizeof(this.C) + binary.Sizeof(int16(this.B))
			return size
		}
		func (this *S) Encode() ([]byte, error) {
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
	var s, ss S
	s.A = 0x11223344
	s.B = -5
	s.C = "hello"

	b, err := binary.Encode(&s, nil)
	if err != nil {
		fmt.Println("binary.Encode failed:", err)
	}

	err = binary.Decode(b, &ss)
	if err != nil {
		fmt.Println("binary.Decode failed:", err)
	}
	fmt.Printf("%+v\n%#v\n%+v", s, b, ss)

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
	size := binary.Sizeof(s)

	fmt.Printf("Sizeof(s)  = %d\nstd.Size(s)= %d\ngob.Size(s)= %d", size, stdSize, gobSize)

	// Output:
	// Sizeof(s)  = 133
	// std.Size(s)= 217
	// gob.Size(s)= 412
}

func ExampleRegStruct() {
	type someRegedStruct struct {
		A int    `binary:"ignore"`
		B uint64 `binary:"packed"`
		C string
		D uint
	}
	binary.RegStruct((*someRegedStruct)(nil))

	// Output:
}

func ExampleRegStruct_packedInts() {
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
	binary.RegStruct((*regedPackedInts)(nil))

	var ints = regedPackedInts{1, 2, 3, 4, 5, 6, []uint64{7, 8, 9}, 10}
	b, err := binary.Encode(ints, nil)
	if err != nil {
		fmt.Println(err)
	}

	if size := binary.Sizeof(ints); size != len(b) {
		fmt.Printf("PackedInts got %+v %+v\nneed %+v\n", len(b), b, size)
	}

	fmt.Printf("Encode packed ints:\n%+v\nsize=%d result=%#v", ints, len(b), b)

	// Output:
	// Encode packed ints:
	// {A:1 B:2 C:3 D:4 E:5 F:6 G:[7 8 9] H:10}
	// size=10 result=[]byte{0x2, 0x4, 0x6, 0x4, 0x5, 0x6, 0x3, 0x7, 0x8, 0x9}
}
