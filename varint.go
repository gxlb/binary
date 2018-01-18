// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary

// This file implements "varint" encoding of 64-bit integers.
// The encoding is:
// - unsigned integers are serialized 7 bits at a time, starting with the
//   least significant bits
// - the most significant bit (msb) in each output byte indicates if there
//   is a continuation byte (msb = 1)
// - signed integers are mapped to unsigned integers using "zig-zag"
//   encoding: Positive values x are written as 2*x + 0, negative values
//   are written as 2*(^x) + 1; that is, negative numbers are complemented
//   and whether to complement is encoded in bit 0.
//
// Design note:
// At most 10 bytes are needed for 64-bit values. The encoding could
// be more dense: a full 64-bit value needs an extra byte just to hold bit 63.
// Instead, the msb of the previous byte could be used to hold bit 63 since we
// know there can't be more than 64 bits. This is a trivial improvement and
// would reduce the maximum encoding length to 9 bytes. However, it breaks the
// invariant that the msb is always the "continuation bit" and thus makes the
// format incompatible with a varint encoding for larger numbers (say 128-bit).

import (
	"errors"
	"fmt"
	"io"
)

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 9

	longUvarintFlagMask     = byte(0x80)
	shortUvarintMaxByteNum  = 2
	shortUvarintLenBitNum   = 2 //length bits of short uvarint(<= 2 bytes) 00~01
	longUvarintLenBitNum    = 4 //length bits of long uvarint(> 2 bytes) 1000~1111
	shortUvarintValueBitNum = 8 - shortUvarintLenBitNum
	longUvarintValueBitNum  = 8 - longUvarintLenBitNum
	shortUvarintValueMask   = 1<<shortUvarintValueBitNum - 1
	longUvarintValueMask    = 1<<longUvarintValueBitNum - 1
)

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
//func PutUvarint(buf []byte, x uint64) int {
//	headByte, followByteNum := packUvarintHead(x)
//	buf[0] = headByte
//	if followByteNum > 0 {
//		for i, x_ := uint8(1), x; i <= followByteNum; i++ {
//			buf[i] = byte(x_)
//			x_ >>= 8
//		}
//	}
//	return int(followByteNum + 1)
//}

// Uvarint decodes a uint64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read
//
func Uvarint(buf []byte) (uint64, int) {
	headByte := buf[0]
	followByteNum, topBits := unpackUvarintHead(headByte)
	size := int(followByteNum + 1)
	if size > MaxVarintLen64 {
		panic(fmt.Errorf("binary.Uvarint: overflow 64-bits value len=%d", size))
	}
	x := topBits
	if followByteNum > 0 {
		for i, shift := uint8(1), uint(0); i <= followByteNum; i, shift = i+1, shift+8 {
			x |= uint64(buf[i]) << shift
		}
	}
	return x, size
}

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int {
	return PutUvarint(buf, ToUvarint(x))
}

// Varint decodes an int64 from buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 with the following meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//              and -n is the number of bytes read
//
func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf) // ok to continue in presence of error
	return ToVarint(ux), n
}

var errOverflow = errors.New("binary: varint overflows a 64-bit integer")

// ReadUvarint reads an encoded unsigned integer from r and returns it as a uint64.
func ReadUvarint(r io.Reader) (uint64, error) {
	var buff [8]byte
	if _, err := r.Read(buff[:1]); err != nil {
		return 0, err
	}
	headByte := buff[0]
	followByteNum, topBits := unpackUvarintHead(headByte)
	if size := int(followByteNum + 1); size > MaxVarintLen64 {
		//return 0, 0
		return 0, fmt.Errorf("binary.ReadUvarint: overflow 64-bits value len=%d ", size)
	}
	x := topBits
	if followByteNum > 0 {
		if n, err := r.Read(buff[:followByteNum]); n != int(followByteNum) || err != nil {
			return 0, err
		}
		for i, shift := uint8(0), uint(0); i < followByteNum; i, shift = i+1, shift+8 {
			x |= uint64(buff[i]) << shift
		}
	}
	return x, nil
}

// ReadVarint reads an encoded signed integer from r and returns it as an int64.
func ReadVarint(r io.Reader) (int64, error) {
	ux, err := ReadUvarint(r) // ok to continue in presence of error
	return ToVarint(ux), err
}

////////////////////////////////////////////////////////////////////////////////

// ToUvarint convert an int64 value to uint64 ZigZag-encoding value for encoding.
// Different from uint64(x), it will move sign bit to bit 0.
// To help to cost fewer bytes for little negative numbers.
// eg: -5 will be encoded as 0x9.
func ToUvarint(x int64) uint64 {
	ux := uint64(x) << 1 // move sign bit to bit0
	if x < 0 {
		ux = ^ux
	}
	return ux
}

// ToVarint decode an uint64 ZigZag-encoding value to original int64 value.
func ToVarint(ux uint64) int64 {
	x := int64(ux >> 1) //move bit0 to sign bit
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

// SizeofVarint return bytes number of an int64 value store as varint
func SizeofVarint(x int64) int {
	return SizeofUvarint(ToUvarint(x))
}

// SizeofUvarint return bytes number of an uint64 value store as uvarint
func SizeofUvarint(ux uint64) int {
	//	i := 1
	//	for n := ux; n >= 0x40; n >>= 8 { //short style
	//		i++
	//	}
	//	if i >= longUvarintMinBytes && n >= 0x10 { //long style, check if need more bytes
	//		i++
	//	}
	//	return i
	size, _ := sizeofUvarint(ux)
	return size
}

func PutUvarint(buf []byte, ux uint64) (size int) {
	n, x := 1, ux
	for ; x > 0x3F; x >>= 8 { //short style, 6 effective bits
		buf[n] = byte(x)
		n++
	}
	headByte := byte(x)
	followByteNum := uint8(n - 1)
	if n > shortUvarintMaxByteNum { //long style, 4 effective bits, check if need more bytes
		if x > 0x0F {
			buf[n] = byte(x)
			n, headByte = n+1, 0
			followByteNum++
		}
		headByte |= longUvarintFlagMask | ((followByteNum - shortUvarintMaxByteNum) << longUvarintValueBitNum)
	} else { //short style, 6 effective bits
		headByte |= followByteNum << shortUvarintValueBitNum
	}
	buf[0] = headByte
	return n
}

// SizeofUvarint return bytes number of an uint64 value store as uvarint
func sizeofUvarint(ux uint64) (size int, topBits byte) {
	n, x := 1, ux
	for ; x > 0x3F; x >>= 8 { //short style, 6 effective bits
		n++
	}
	if n > shortUvarintMaxByteNum && x > 0x0F { //long style, 4 effective bits, check if need more bytes
		n, x = n+1, 0
	}
	return n, byte(x)
}

// 00~01 1~2 bytes 0~14 bits
// 100~101 3~4
// 1100~1101 5~6
// 11100~11101 7~10

// 00~01 1~2 bytes 0~14 bits 0~16383
// 1000~1111 3~9 bytes
func packUvarintHead(ux uint64) (headByte byte, followByteNum uint8) {
	size, topBits := sizeofUvarint(ux)
	followByteNum = byte(size - 1)
	if size <= shortUvarintMaxByteNum { //short style
		headByte = followByteNum << shortUvarintValueBitNum
	} else { //long style
		headByte = longUvarintFlagMask | ((followByteNum - shortUvarintMaxByteNum) << longUvarintValueBitNum)
	}
	headByte |= topBits
	return
}

func unpackUvarintHead(headByte byte) (followByteNum uint8, topBits uint64) {
	if headByte&longUvarintFlagMask == 0 { //short style
		followByteNum = headByte >> shortUvarintValueBitNum
		topBits = uint64(headByte & shortUvarintValueMask)
	} else { //long style
		followByteNum = (headByte&0x7f)>>longUvarintValueBitNum + shortUvarintMaxByteNum
		topBits = uint64(headByte & shortUvarintValueMask)
	}
	topBits <<= (8 * followByteNum)
	return
}
