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
// Package binary is uesed to Encode/Decode between go data and byte slice.
//
// The main purpose of this package is to replace package "std.binary".
// The design goal is to take both advantages of std.binary(encoding/binary) and gob.
//
// Upgraded from std.binary(encoding/binary).
//
// Compare with other serialization package, this package is with full-feature as
// gob and protocol buffers, and with high-performance and lightweight as std.binary.
// It is designed as a common solution to easily encode/decode between go data and byte slice.
// It is recommended to use in net protocol serialization and go memory data serialization such as DB.
//
// Support all serialize-able data types:
//
//	int, int8, int16, int32, int64,
//	uint, uint8, uint16, uint32, uint64,
//	float32, float64, complex64, complex128,
//	bool, string, slice, array, map, struct.
//	int/uint will be encoded as varint(1~10 bytes).
//	And their direct pointers.
//	eg: *string, *struct, *map, *slice, *int32.
//
// Here is the main feature of this package.
//  1. as light-weight as std.binary
//  2. with full-type support like gob.
//  3. as high-performance as std.binary and gob.
//  4. encoding with fewer bytes than std.binary and gob.
//  5. use RegStruct to improve performance of struct encoding/decoding
//  6. recommended using in net protocol serialization and DB serialization
//
// Under MIT license.
//
//	Copyright (c) 2017 Ally Dale<vipally@gmail.com>
//	Author  : Ally Dale<vipally@gmail.com>
//	Site    : https://github.com/vipally
//	Origin  : https://github.com/vipally/binary
package binary

// **************************************************************************
// TODO:
// 1.[Encoder/Decoder].RegStruct
// 2.[Encoder/Decoder].RegMarshaler
// 3.[Encoder/Decoder].ResizeBuffer
// 4.bool as a bit
// 5.pointer put a bool to check if it is nil
// 6.reg interface to using PackUnpacker interface
// 7.use `binary:"packed"` to set if store a int_n value as varint
// ***************************************************************************

import (
	"errors"
	"io"
	"reflect"
)

// Size is same to Sizeof.
// Size returns how many bytes Write would generate to encode the value v, which
// must be a serialize-able value or a slice/map of serialize-able values, or a pointer to such data.
// If v is neither of these, Size returns -1.
func Size(data interface{}) int {
	return Sizeof(data)
}

// Sizeof returns how many bytes Write would generate to encode the value v, which
// must be a serialize-able value or a slice/map/struct of serialize-able values, or a pointer to such data.
// If v is neither of these, Size returns -1.
// If data implements interface BinarySizer, it will use data.Size first.
// It will panic if data implements interface BinarySizer or BinaryEncoder only.
func Sizeof(data interface{}) int {
	if p, ok := data.(BinarySizer); ok {
		if _, _ok := data.(BinaryEncoder); !_ok { //interface verification
			panic(errors.New("expect but not BinaryEncoder:" + reflect.TypeOf(data).String()))
		}
		return p.Size()
	} else {
		if _, _ok := data.(BinaryEncoder); _ok { //interface verification
			panic(errors.New("unexpected BinaryEncoder:" + reflect.TypeOf(data).String()))
		}
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
	var decoder Decoder
	decoder.Init(nil, endian)
	decoder.reader = r
	return decoder.Value(data)
}

// Write writes the binary representation of data into w.
// Data must be a fixed-size value or a slice of fixed-size
// values, or a pointer to such data.
// Boolean values encode as one byte: 1 for true, and 0 for false.
// Bytes written to w are encoded using the specified byte order
// and read from successive fields of the data.
// When writing structs, zero values are written for fields
// with blank (_) field names.
func Write(w io.Writer, endian Endian, data interface{}) error {
	size := Sizeof(data)
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

	err := encoder.Value(data)
	assert(err == nil, err) //Value will never return error, because Sizeof(data) < 0 has blocked the error data

	_, err = w.Write(encoder.Buffer())
	return err
}

// BinarySizer is an interface to define go data Size method.
type BinarySizer interface {
	Size() int
}

// BinaryEncoder is an interface to define go data Encode method.
// buffer is nil-able.
type BinaryEncoder interface {
	Encode(buffer []byte) ([]byte, error)
}

// BinaryDecoder is an interface to define go data Decode method.
type BinaryDecoder interface {
	Decode(buffer []byte) error
}

// interface BinarySerializer defines the go data Size/Encode/Decode method
type BinarySerializer interface {
	BinarySizer
	BinaryEncoder
	BinaryDecoder
}

// Encode marshal go data to byte array.
// nil buffer is aviable, it will create new buffer if necessary.
func Encode(data interface{}, buffer []byte) ([]byte, error) {
	buff, err := MakeEncodeBuffer(data, buffer)
	if err != nil {
		return nil, err
	}

	encoder := NewEncoderBuffer(buff)

	err = encoder.Value(data)
	return encoder.Buffer(), err
}

// Decode unmarshal go data from byte array.
// data must be interface of pointer for modify.
// It will make new pointer or slice/map for nil-field of data.
func Decode(buffer []byte, data interface{}) error {
	var decoder Decoder
	decoder.Init(buffer, DefaultEndian)
	return decoder.Value(data)
}

// MakeEncodeBuffer create enough buffer to encode data.
// nil buffer is aviable, it will create new buffer if necessary.
func MakeEncodeBuffer(data interface{}, buffer []byte) ([]byte, error) {
	size := Sizeof(data)
	if size < 0 {
		return nil, errors.New("binary.MakeEncodeBuffer: invalid type " + reflect.TypeOf(data).String())
	}

	buff := buffer
	if len(buff) < size {
		buff = make([]byte, size)
	}
	return buff, nil
}
