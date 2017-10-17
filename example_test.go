// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary_test

import (
	"bytes"
	"fmt"
	"math"

	"github.com/vipally/binary"
)

func ExampleWrite() {
	buf := new(bytes.Buffer)
	var pi float64 = math.Pi
	err := binary.Write(buf, binary.LittleEndian, pi)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	fmt.Printf("% x", buf.Bytes())
	// Output: 18 2d 44 54 fb 21 09 40
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
	fmt.Printf("%x", buf.Bytes())
	// Output: beefcafe
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
	// Output: 3.141592653589793
}

func ExamplePack() {
	var s struct {
		A uint32
		B int
		C string
	}
	s.A = 0x11223344
	s.B = -5
	s.C = "hello"
	b, err := binary.Pack(s, nil)
	if err != nil {
		fmt.Println("binary.Pack failed:", err)
	}
	fmt.Printf("%#v", b)
	// Output: []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExamplePack_withbuffer() {
	var s struct {
		A uint32
		B int
		C string
	}
	s.A = 0x11223344
	s.B = -5
	s.C = "hello"
	size := binary.Sizeof(s)
	buffer := make([]byte, size)
	b, err := binary.Pack(s, buffer)
	if err != nil {
		fmt.Println("binary.Pack failed:", err)
	}
	fmt.Printf("%#v", b)
	// Output: []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleUnack() {
	var s struct {
		A uint32
		B int
		C string
	}
	buffer := []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	err := binary.Unpack(buffer, &s)
	if err != nil {
		fmt.Println("binary.Unpack failed:", err)
	}
	fmt.Printf("%#v", s)
	// Output: struct { A uint32; B int; C string }{A:0x11223344, B:-5, C:"hello"}
}
func ExampleEncoder() {
	encoder := binary.NewEncoder(100)
	encoder.Uint32(0x11223344)
	encoder.Varint(-5)
	encoder.String("hello")
	encodeResult := encoder.Buffer()
	fmt.Printf("%#v", encodeResult)
	// Output: []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
}
func ExampleDecoder() {
	buffer := []byte{0x44, 0x33, 0x22, 0x11, 0x9, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	decoder := binary.NewDecoder(buffer)
	u32 := decoder.Uint32()
	i, _ := decoder.Varint()
	str := decoder.String()
	fmt.Printf("[%#v %#v %#v]", u32, i, str)
	// Output: [0x11223344 -5 "hello"]
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
func (this *S) Pack() ([]byte, error) {
	encoder := binary.NewEncoder(this.Size())
	encoder.Value(this.A)
	encoder.Int16(int16(this.B))
	encoder.Value(this.C)
	return encoder.Buffer(), nil
}
func (this *S) Unpack(buffer []byte) error {
	decoder := binary.NewDecoder(buffer)
	decoder.Value(&this.A)
	this.B = int(decoder.Int16())
	decoder.Value(&this.C)
	return nil
}
func ExamplePacker() {
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
		func (this *S) Pack() ([]byte, error) {
			encoder := binary.NewEncoder(this.Size())
			encoder.Value(this.A)
			encoder.Int16(int16(this.B))
			encoder.Value(this.C)
			return encoder.Buffer(), nil
		}
		func (this *S) Unpack(buffer []byte) error {
			decoder := binary.NewDecoder(buffer)
			decoder.Value(&this.A)
			this.B = int(decoder.Int16())
			decoder.Value(&this.C)
			return nil
		}
	*/
	var s, ss S
	s.A = 0x11223344
	s.B = -5
	s.C = "hello"
	b, err := binary.Pack(&s, nil)

	if err != nil {
		fmt.Println("binary.Pack failed:", err)
	}
	err = binary.Unpack(b, &ss)
	if err != nil {
		fmt.Println("binary.Unpack failed:", err)
	}
	fmt.Printf("[%v\n%#v\n%v]", s, b, ss)
	// Output: [{287454020 -5 hello}
	//[]byte{0x44, 0x33, 0x22, 0x11, 0xfb, 0xff, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	//{287454020 -5 hello}]
}
