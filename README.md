# binary [![Build Status](https://travis-ci.org/vipally/binary.svg?branch=master)](https://travis-ci.org/vipally/binary) [![Coverage Status](https://coveralls.io/repos/vipally/binary/badge.svg?branch=master)](https://coveralls.io/r/vipally/binary?branch=master) [![GoDoc](https://godoc.org/github.com/vipally/binary?status.svg)](https://godoc.org/github.com/vipally/binary) ![Version](https://img.shields.io/badge/version-0.8.0.final-green.svg)
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

# 2. [recommended usage] Use Pack/UnPack to read/write memory buffer directly.
	If data implement interface Packer, it will use data.Pack/data.Unpack 
	to encode/decode data.
	NOTE that data.Unpack must implement on pointer receiever to enable modifying
	receiever.Even though Size/Pack of data can implement on non-pointer receiever,
	binary.Pack(&data, nil) is required if data has implement interface Packer.
	binary.Pack(data, nil) will probably NEVER use Packer methods to Pack/Unpack
	data.
	eg:

	import "github.com/vipally/binary"
	
	//1.Pack with default buffer
	if bytes, err := binary.Pack(&data, nil); err==nil{
		//...
	}

	//2.Pack with existing buffer
	size := binary.Sizeof(data)
	buffer := make([]byte, size)
	if bytes, err := binary.Pack(&data, buffer); err==nil{
		//...
	}

	//3.Unpack from buffer
	if err := binary.Unpack(bytes, &data); err==nil{
		//...
	}

# 3. [advanced usage] Encoder/Decoder are exported types aviable for encoding/decoding.
	eg:
	encoder := binary.NewEncoder(bufferSize)
	encoder.Uint32(u32)
	encoder.String(str)
	encodeResult := encoder.Buffer()
	
	decoder := binary.NewDecoder(buffer)
	u32 := decoder.Uint32()
	str := decoder.String()

# 4. Put an extra length field(uvarint,1~10 bytes) before string, slice, array, map.
	eg: 
	var s string = "hello"
	will be encoded as:
	[]byte{0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

# 5. Pack bool array with bits.
	eg: 
	[]bool{true, true, true, false, true, true, false, false, true}
	will be encoded as:
	[]byte{0x9, 0x37, 0x1}

# 6. Hide struct field when encoding/decoding.
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

# 7. Auto allocate for slice, map and pointer.
	eg: 
	type S struct{
	    A *uint32
		B *string
		C *[]uint8
		D []uint32
	}
	It will new pointers for fields "A, B, C",
	and make new slice for fields "*C, D" when decode.
	
# 8. int/uint values will be encoded as varint/uvarint(1~10 bytes).
	eg: 
	uint(1)     will be encoded as: []byte{0x1}
	uint(128)   will be encoded as: []byte{0x80, 0x1}
	uint(32765) will be encoded as: []byte{0xfd, 0xff, 0x1}
	int(-5)     will be encoded as: []byte{0x9}
	int(-65)    will be encoded as: []byte{0x81, 0x1}
	
	

