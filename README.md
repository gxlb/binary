# binary [![GoDoc](https://godoc.org/github.com/vipally/binary?status.svg)](https://godoc.org/github.com/vipally/binary) ![Version](https://img.shields.io/badge/version-0.8.0.final-green.svg)
  Package binary is uesed to Pack/Unpack between go data and byte slice.

  The main purpose of this package is to replace package "std.binary".

***

CopyRight 2017 @Ally Dale. All rights reserved.
	
Author  : [Ally Dale(vipally@gmail.com)](mailto://vipally@gmail.com)

Blog    : [http://blog.csdn.net/vipally](http://blog.csdn.net/vipally)

Site    : [https://github.com/vipally](https://github.com/vipally)

****

# 1. support all serialize-able basic types:
	int, int8, int16, int32, int64,
	uint, uint8, uint16, uint32, uint64,
	float32, float64, complex64, complex128,
	bool, string, slice, array, map, struct.
	And their direct pointers. eg: *string, *struct, *map, *slice, *int32.

# 2. store an extra length field(uvarint,1~10 bytes) for string, slice, array, map.
	eg: 
	var s string = "hello"
	will be encode as:
	[]byte{0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

# 3. pack bool array with bits.
	eg: 
	[]bool{true, true, true, false}
	will be encoded as:
	[]byte{0x4, 0x7}

# 4. hide struct field when encode/decode.
	Only encode/decode exported fields.
	Support field tag `binary:"ignore"` to disable encode/decode fields.
	eg: 
	type S struct{
	    A uint32
		b uint32
		_ uint32
		C uint32 `binary:"ignore"`
	}
	only field "A" will be encode/decode.

# 5. auto alloc for slice, map and pointer.
	eg: 
	type S struct{
	    A *uint32
		B *string
		C *[]uint8
		D []uint32
	}
	It will new pointers for fields "A, B, C",
	and make new slice for fields "*C, D" when decode.
