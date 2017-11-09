package binary

import (
	"errors"
	"fmt"
	"io"
)

var (
	// ErrNotEnoughSpace buffer not enough
	ErrNotEnoughSpace = errors.New("not enough space")
)

type coder struct {
	buff []byte
	pos  int

	boolPos int  //index of last bool set in buff
	boolBit byte //bit of next aviable bool
	endian  Endian
}

func (cder *coder) setEndian(endian Endian) {
	cder.endian = endian
}

// Buffer returns the byte slice that has been encoding/decoding.
func (cder *coder) Buffer() []byte {
	return cder.buff[:cder.pos]
}

// Len returns unmber of bytes that has been encoding/decoding.
func (cder *coder) Len() int {
	return cder.pos
}

// Cap returns number total bytes of cder coder buffer.
func (cder *coder) Cap() int {
	return len(cder.buff)
}

// Skip ignore the next size of bytes for encoding/decoding and
// set skiped bytes to 0.
// It will panic if space not enough.
// It will return -1 if size <= 0.
func (cder *coder) Skip(size int) int {
	newPos := cder.pos + size
	if size >= 0 && newPos <= cder.Cap() {
		for i, b := int(size-1), cder.buff[cder.pos:newPos]; i >= 0; i-- { //zero skiped bytes
			b[i] = 0
		}
		cder.pos = newPos
		return size
	}
	return -1
}

// Reset move the read/write pointer to the beginning of buffer
// and set all reseted bytes to 0.
func (cder *coder) Reset() {
	for i := int(cder.pos - 1); i >= 0; i-- { //zero encoded bytes
		cder.buff[i] = 0
	}
	cder.pos = 0
	cder.resetBoolCoder()
}

// reset the state of bool coder
func (cder *coder) resetBoolCoder() {
	cder.boolPos = -1
	cder.boolBit = 0
}

// reserve returns next size bytes for encoding/decoding.
// it will panic if not enough space.
func (cder *coder) reserve(size int) []byte {
	newPos := cder.pos + size
	if newPos > cder.Cap() {
		panic(fmt.Errorf("binary.Coder:buffer overflow pos=%d cap=%d require=%d, not enough space", cder.pos, cder.Cap(), size))
	}
	if size > 0 && newPos <= cder.Cap() {
		b := cder.buff[cder.pos:newPos]
		cder.pos = newPos
		return b
	}
	return nil
}

// BytesReader transform bytes as Reader
type BytesReader []byte

// Read from bytes
func (p *BytesReader) Read(data []byte) (n int, err error) {
	n = copy(data, *p)
	if n == len(*p) {
		err = io.EOF
	}
	*p = (*p)[n:]
	return
}

// BytesWriter transform bytes as Writer
type BytesWriter []byte

// Write to bytes
func (p *BytesWriter) Write(data []byte) (n int, err error) {
	n = copy(*p, data)
	if n < len(data) {
		err = ErrNotEnoughSpace
	}
	*p = (*p)[n:]
	return
}
