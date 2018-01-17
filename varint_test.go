// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package binary

import (
	"bytes"
	std "encoding/binary"
	"fmt"
	"io"
	"testing"
)

var (
	UvarintCases = []uint64{
		0x0000000000000001, 0x0000000000000003, 0x0000000000000007, 0x000000000000000F,
		0x000000000000001F, 0x000000000000003F, 0x000000000000007F, 0x00000000000000FF,
		0x00000000000001FF, 0x00000000000003FF, 0x00000000000007FF, 0x0000000000000FFF,
		0x0000000000001FFF, 0x0000000000003FFF, 0x0000000000007FFF, 0x000000000000FFFF,
		0x000000000001FFFF, 0x000000000003FFFF, 0x000000000007FFFF, 0x00000000000FFFFF,
		0x00000000001FFFFF, 0x00000000003FFFFF, 0x00000000007FFFFF, 0x0000000000FFFFFF,
		0x0000000001FFFFFF, 0x0000000003FFFFFF, 0x0000000007FFFFFF, 0x000000000FFFFFFF,
		0x000000001FFFFFFF, 0x000000003FFFFFFF, 0x000000007FFFFFFF, 0x00000000FFFFFFFF,
		0x00000001FFFFFFFF, 0x00000003FFFFFFFF, 0x00000007FFFFFFFF, 0x0000000FFFFFFFFF,
		0x0000001FFFFFFFFF, 0x0000003FFFFFFFFF, 0x0000007FFFFFFFFF, 0x000000FFFFFFFFFF,
		0x000001FFFFFFFFFF, 0x000003FFFFFFFFFF, 0x000007FFFFFFFFFF, 0x00000FFFFFFFFFFF,
		0x00001FFFFFFFFFFF, 0x00003FFFFFFFFFFF, 0x00007FFFFFFFFFFF, 0x0000FFFFFFFFFFFF,
		0x0001FFFFFFFFFFFF, 0x0003FFFFFFFFFFFF, 0x0007FFFFFFFFFFFF, 0x000FFFFFFFFFFFFF,
		0x001FFFFFFFFFFFFF, 0x003FFFFFFFFFFFFF, 0x007FFFFFFFFFFFFF, 0x00FFFFFFFFFFFFFF,
		0x01FFFFFFFFFFFFFF, 0x03FFFFFFFFFFFFFF, 0x07FFFFFFFFFFFFFF, 0x0FFFFFFFFFFFFFFF,
		0x1FFFFFFFFFFFFFFF, 0x3FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
	}
	UvarintStdBytes = []byte{}
	UvarintBytes    = []byte{}
)

func init() {
	//	for i := uint(1); i <= 64; i++ {
	//		x := uint64(1<<i) - 1
	//		fmt.Printf("0x%016X, ", x)
	//		if i%4 == 0 {
	//			fmt.Println("")
	//		}
	//	}
	var buf [10]byte
	cnt := 0

	cnt = 0
	for _, x := range UvarintCases {
		n := std.PutUvarint(buf[0:], x)
		for i := 0; i < n; i++ {
			fmt.Printf("0x%02x, ", buf[i])
			cnt++
			if cnt%16 == 0 {
				fmt.Println("")
			}
		}
	}
	fmt.Println("\ncnt=", cnt)

	cnt = 0
	for _, x := range UvarintCases {
		n := PutUvarint(buf[0:], x)
		for i := 0; i < n; i++ {
			fmt.Printf("0x%02x, ", buf[i])
			cnt++
			if cnt%16 == 0 {
				fmt.Println("")
			}
		}
	}
	fmt.Println("\ncnt=", cnt)
}

func testConstant(t *testing.T, w uint, max int) {
	buf := make([]byte, MaxVarintLen64)
	n := PutUvarint(buf, 1<<w-1)
	if n != max {
		t.Errorf("MaxVarintLen%d = %d; want %d", w, max, n)
	}
}

func TestConstants(t *testing.T) {
	testConstant(t, 16, MaxVarintLen16)
	testConstant(t, 32, MaxVarintLen32)
	testConstant(t, 64, MaxVarintLen64)
}

func testVarint(t *testing.T, x int64) {
	buf := make([]byte, MaxVarintLen64)
	n := PutVarint(buf, x)
	y, m := Varint(buf[0:n])
	if x != y {
		t.Errorf("Varint(%d): got %d", x, y)
	}
	if n != m {
		t.Errorf("Varint(%d): got n = %d; want %d", x, m, n)
	}

	y, err := ReadVarint(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("ReadVarint(%d): %s", x, err)
	}
	if x != y {
		t.Errorf("ReadVarint(%d): got %d", x, y)
	}
}

func testUvarint(t *testing.T, x uint64) {
	buf := make([]byte, MaxVarintLen64)
	n := PutUvarint(buf, x)
	y, m := Uvarint(buf[0:n])
	if x != y {
		t.Errorf("Uvarint(%d): got %d", x, y)
	}
	if n != m {
		t.Errorf("Uvarint(%d): got n = %d; want %d", x, m, n)
	}

	y, err := ReadUvarint(bytes.NewReader(buf))
	if err != nil {
		t.Errorf("ReadUvarint(%d): %s", x, err)
	}
	if x != y {
		t.Errorf("ReadUvarint(%d): got %d", x, y)
	}
}

var tests = []int64{
	-1 << 63,
	-1<<63 + 1,
	-1,
	0,
	1,
	2,
	10,
	20,
	63,
	64,
	65,
	127,
	128,
	129,
	255,
	256,
	257,
	1<<63 - 1,
}

func TestVarint(t *testing.T) {
	for _, x := range tests {
		testVarint(t, x)
		testVarint(t, -x)
	}
	for x := int64(0x7); x != 0; x <<= 1 {
		testVarint(t, x)
		testVarint(t, -x)
	}
}

func TestUvarint(t *testing.T) {
	for _, x := range tests {
		testUvarint(t, uint64(x))
	}
	for x := uint64(0x7); x != 0; x <<= 1 {
		testUvarint(t, x)
	}
}

func TestBufferTooSmall(t *testing.T) {
	buffer := []byte{0x80, 0x80, 0x80, 0x80}
	for i := 0; i <= len(buffer); i++ {
		buf := buffer[0:i]
		x, n := Uvarint(buf)
		if x != 0 || n != 0 {
			t.Errorf("Uvarint(%v): got x = %d, n = %d", buf, x, n)
		}

		x, err := ReadUvarint(bytes.NewReader(buf))
		if x != 0 || err != io.EOF {
			t.Errorf("ReadUvarint(%v): got x = %d, err = %s", buf, x, err)
		}
	}
}

func testOverflow(t *testing.T, buf []byte, n0 int, err0 error) {
	x, n := Uvarint(buf)
	if x != 0 || n != n0 {
		t.Errorf("Uvarint(%v): got x = %d, n = %d; want 0, %d", buf, x, n, n0)
	}

	x, err := ReadUvarint(bytes.NewReader(buf))
	if x != 0 || err != err0 {
		t.Errorf("ReadUvarint(%v): got x = %d, err = %s; want 0, %s", buf, x, err, err0)
	}
}

func TestOverflow(t *testing.T) {
	testOverflow(t, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x2}, -10, errOverflow)
	testOverflow(t, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x1, 0, 0}, -13, errOverflow)
}

func TestNonCanonicalZero(t *testing.T) {
	buf := []byte{0x80, 0x80, 0x80, 0}
	x, n := Uvarint(buf)
	if x != 0 || n != 4 {
		t.Errorf("Uvarint(%v): got x = %d, n = %d; want 0, 4", buf, x, n)

	}
}

func BenchmarkPutUvarint32(b *testing.B) {
	buf := make([]byte, MaxVarintLen32)
	b.SetBytes(4)
	for i := 0; i < b.N; i++ {
		for j := uint(0); j < MaxVarintLen32; j++ {
			PutUvarint(buf, 1<<(j*7))
		}
	}
}

func BenchmarkPutUvarint64(b *testing.B) {
	buf := make([]byte, MaxVarintLen64)
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		for j := uint(0); j < MaxVarintLen64; j++ {
			PutUvarint(buf, 1<<(j*7))
		}
	}
}
