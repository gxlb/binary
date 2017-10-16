// Package binary implements simple translation between numbers and byte
// sequences and encoding and decoding of varints.
//
// Numbers are translated by reading and writing fixed-size values.
// A fixed-size value is either a fixed-size arithmetic
// type (bool, int8, uint8, int16, float32, complex64, ...)
// or an array or struct containing only fixed-size values.
//
// The varint functions encode and decode single integer values using
// a variable-length encoding; smaller values require fewer bytes.
// For a specification, see
// https://developers.google.com/protocol-buffers/docs/encoding.
//
// This package favors simplicity over efficiency. Clients that require
// high-performance serialization, especially for large data structures,
// should look at more advanced solutions such as the encoding/gob
// package or protocol buffers.
//
//
// Author: Ally Dale<vipally@gmail.com>
//
// Package binary is uesed to Pack/Unpack between go data and byte slice.
//
// The main purpose of this package is to replace package "std.binary".
//
// Support all serialize-able data types:
//
//    bool, int8, int16, int32, int64,
//    uint8, uint16, uint32, uint64,
//    float32, float64, complex64, complex128,
//    string, struct, slice, array, map.
//    And their direct pointers.
package binary

import (
	"errors"
	"io"
	"reflect"
)

// Size returns how many bytes Write would generate to encode the value v, which
// must be a fixed-size value or a slice of fixed-size values, or a pointer to such data.
// If v is neither of these, Size returns -1.
func Size(data interface{}) int {
	if p, ok := data.(Packer); ok {
		return p.Size()
	}
	return sizeof(data)
}

// Read reads structured binary data from r into data.
// Data must be a pointer to a fixed-size value or a slice
// of fixed-size values.
// Bytes read from r are decoded using the specified byte order
// and written to successive fields of the data.
// When decoding boolean values, a zero byte is decoded as false, and
// any other non-zero byte is decoded as true.
// When reading into structs, the field data for fields with
// blank (_) field names is skipped; i.e., blank field names
// may be used for padding.
// When reading into a struct, all non-blank fields must be exported.
//
// The error is EOF only if no bytes were read.
// If an EOF happens after reading some but not all the bytes,
// Read returns ErrUnexpectedEOF.
func Read(r io.Reader, endian Endian, data interface{}) error {
	//	size := Size(data)
	//	if size < 0 {
	//		return errors.New("binary.Read: invalid type " + reflect.TypeOf(data).String())
	//	}
	//	//	var b [16]byte
	//	//	var bs []byte
	//	//	if size > len(b) {
	//	//		bs = make([]byte, size)
	//	//	} else {
	//	//		bs = b[:size]
	//	//	}
	//	//	if _, err := io.ReadFull(r, bs); err != nil {
	//	//		return err
	//	//	}
	//	b, _ := readAll(r, size)
	var decoder decoderReader
	decoder.Init(r, endian)
	return decoder.Value(data)
}

//func readAll(r io.Reader, size int) ([]byte, error) {
//	var ret []byte = make([]byte, 0, size)
//	var buff [512]byte
//	for {
//		if n, err := r.Read(buff[0:]); n > 0 {
//			ret = append(ret, buff[:n]...)
//			if err == io.EOF {
//				break
//			}
//		} else {
//			break
//		}
//	}
//	return ret, nil
//}

// Write writes the binary representation of data into w.
// Data must be a fixed-size value or a slice of fixed-size
// values, or a pointer to such data.
// Boolean values encode as one byte: 1 for true, and 0 for false.
// Bytes written to w are encoded using the specified byte order
// and read from successive fields of the data.
// When writing structs, zero values are written for fields
// with blank (_) field names.
func Write(w io.Writer, endian Endian, data interface{}) error {
	size := Size(data)
	if size < 0 {
		return errors.New("binary.Write: invalid type " + reflect.TypeOf(data).String())
	}
	var b [16]byte
	var bs []byte
	if size > len(b) {
		bs = make([]byte, size)
	} else {
		bs = b[:size]
	}
	var encoder Encoder
	encoder.buff = bs
	encoder.endian = endian
	encoder.pos = 0
	if err := encoder.Value(data); err != nil {
		return err
	}
	if _, err := w.Write(encoder.Buffer()); err != nil {
		return err
	}
	return nil
}

// Packer is an interface to define go data Pack and UnPack method.
type Packer interface {
	Sizer
	Pack(buffer []byte) ([]byte, error)
	Unpack(buffer []byte) error
}

// Pack encode go data to byte array.
// Buffer is nil-aviable, it will create new buffer if necessary.
func Pack(data interface{}, buffer []byte) ([]byte, error) {
	size := Size(data)
	if size < 0 {
		return nil, errors.New("binary.Pack: invalid type " + reflect.TypeOf(data).String())
	}
	if len(buffer) < size {
		buffer = make([]byte, size)
	}
	if p, ok := data.(Packer); ok {
		return p.Pack(buffer)
	}
	var encoder Encoder
	encoder.buff = buffer
	encoder.endian = DefaultEndian
	encoder.pos = 0
	if err := encoder.Value(data); err != nil {
		return nil, err
	}
	return encoder.Buffer(), nil
}

// Unpack decode go data from byte array.
func Unpack(buffer []byte, data interface{}) error {
	if p, ok := data.(Packer); ok {
		return p.Unpack(buffer)
	}

	var decoder Decoder
	decoder.Init(buffer, DefaultEndian)
	return decoder.Value(data)
}
