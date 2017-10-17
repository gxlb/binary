// TODO:
// bench test with std.binary and gob
// function test
// field tag parse
// read buffer not enough, need return errr, not panic

package binary

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrNotEnoughSpace = errors.New("not enough space")
)

type coder struct {
	buff   []byte
	pos    int
	endian Endian
}

// Buffer returns the byte slice that has been encoding/decoding.
func (this *coder) Buffer() []byte {
	return this.buff[:this.pos]
}

// Len returns unmber of bytes that has been encoding/decoding.
func (this *coder) Len() int {
	return this.pos
}

// Cap returns number total bytes of this coder buffer.
func (this *coder) Cap() int {
	return len(this.buff)
}

// Skip ignore size bytes for encoding/decoding.
// If with errors, it will return -1
func (this *coder) Skip(size int) int {
	if size >= 0 && this.pos+size <= this.Cap() {
		this.pos += size
		return size
	}
	return -1
}

// Reset move the read/write pointer to the beginning of buffer.
func (this *coder) Reset() {
	this.pos = 0
}

// reserve returns next size bytes for encoding/decoding.
func (this *coder) reserve(size int) []byte {
	if this.pos+size > this.Cap() {
		panic(fmt.Sprintf("Coder:buff overflow pos=%d size=%d, cap=%d\n", this.pos, size, this.Cap()))
	}
	b := this.buff[this.pos : this.pos+size]
	if this.Skip(size) >= 0 {
		return b
	}
	return nil
}

type BytesReader []byte

func (p *BytesReader) Read(data []byte) (n int, err error) {
	n = copy(data, *p)
	if n == len(*p) {
		err = io.EOF
	}
	*p = (*p)[n:]
	return
}

type BytesWriter []byte

func (p *BytesWriter) Write(data []byte) (n int, err error) {
	n = copy(*p, data)
	if n < len(data) {
		err = ErrNotEnoughSpace
	}
	*p = (*p)[n:]
	return
}
