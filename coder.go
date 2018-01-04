package binary

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"unicode"
	"unicode/utf8"
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
		for i, b := size-1, cder.buff[cder.pos:newPos]; i >= 0; i-- { //zero skiped bytes
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
	for i := cder.pos - 1; i >= 0; i-- { //zero encoded bytes
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

// mustReserve returns next size bytes for encoding/decoding.
// it will panic if not enough space.
func (cder *coder) mustReserve(size int) []byte {
	newPos := cder.pos + size
	_cap := len(cder.buff)
	if newPos > _cap {
		panic(fmt.Errorf("binary.Coder:buffer overflow pos=%d cap=%d require=%d, not enough space", cder.pos, cder.Cap(), size))
	}
	if size > 0 {
		b := cder.buff[cder.pos:newPos]
		cder.pos = newPos
		return b
	}
	return nil
}

// reserve returns next size bytes for encoding/decoding.
func (cder *coder) reserve(size int) ([]byte, error) {
	newPos := cder.pos + size
	_cap := len(cder.buff)
	if newPos > _cap {
		return nil, fmt.Errorf("binary.Coder:buffer overflow pos=%d cap=%d require=%d, not enough space", cder.pos, cder.Cap(), size)
	}
	if size > 0 {
		b := cder.buff[cder.pos:newPos]
		cder.pos = newPos
		return b, nil
	}
	return nil, nil
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

type valueFunc func(v reflect.Value, depth uint, packed bool, serializer serializerSwitch) error
