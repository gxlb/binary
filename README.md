# binary [![Coverage Status](https://coveralls.io/repos/vipally/binary/badge.svg?branch=master)](https://coveralls.io/r/vipally/binary?branch=master) [![GoDoc](https://godoc.org/github.com/vipally/binary?status.svg)](https://godoc.org/github.com/vipally/binary) ![Version](https://img.shields.io/badge/version-0.8.0.final-green.svg)
  Package binary is uesed to Pack/Unpack between go data and byte slice.

  The main purpose of this package is to replace package "std.binary".

  Compare with other serialization package, this package is with full-feature as
  gob and protocol buffers, and with high-performance and lightweight as std.binary.

  It is designed as a common solution to easily encode/decode between go data and byte slice.

  It is recommended to use in net protocol serialization and go memory data serialization such as DB.

***

CopyRight 2017 @Ally Dale. All rights reserved.
	
Author  : [Ally Dale(vipally@gmail.com)](mailto://vipally@gmail.com)

Blog    : [http://blog.csdn.net/vipally](http://blog.csdn.net/vipally)

Site    : [https://github.com/vipally](https://github.com/vipally)

****

# 1. Support all serialize-able basic types:
	int, int8, int16, int32, int64,
	uint, uint8, uint16, uint32, uint64,
	float32, float64, complex64, complex128,
	bool, string, slice, array, map, struct.
	And their direct pointers. 
	eg: *string, *struct, *map, *slice, *int32.

# 2. Store an extra length field(uvarint,1~10 bytes) for string, slice, array, map.
	eg: 
	var s string = "hello"
	will be encoded as:
	[]byte{0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

# 3. Pack bool array with bits.
	eg: 
	[]bool{true, true, true, false, true, true, false, false, true}
	will be encoded as:
	[]byte{0x9, 0x37, 0x1}

# 4. Hide struct field when encode/decode.
	Only encode/decode exported fields.
	Support field tag `binary:"ignore"` to disable encode/decode fields.
	eg: 
	type S struct{
	    A uint32
		b uint32
		_ uint32
		C uint32 `binary:"ignore"`
	}
	Only field "A" will be encode/decode.

# 5. Auto allocate for slice, map and pointer.
	eg: 
	type S struct{
	    A *uint32
		B *string
		C *[]uint8
		D []uint32
	}
	It will new pointers for fields "A, B, C",
	and make new slice for fields "*C, D" when decode.
	
# 6. Use Pack/UnPack read/write memory buffer directly.
	If data implement interface Packer, it will use data.Pack/data.Unpack 
	to encode/decode data.
	eg:
	
	if bytes, err := binary.Pack(&data, nil); err==nil{
		//...
	}

	size := binary.Sizeof(data)
	buffer := make([]byte, size)
	if bytes, err := binary.Pack(&data, buffer); err==nil{
		//...
	}

	if err := binary.Unpack(bytes, &data); err==nil{
		//...
	}

# 7. Encoder/Decoder are exported types aviable for encoding/decoding.
	eg:
	encoder := binary.NewEncoder(bufferSize)
	encoder.Uint32(u32)
	encoder.String(str)
	encodeResult := encoder.Buffer()
	
	decoder := binary.NewDecoder(buffer)
	u32 := decoder.Uint32()
	str := decoder.String()
