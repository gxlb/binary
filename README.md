# binary [![Build Status](https://travis-ci.org/vipally/binary.svg?branch=master)](https://travis-ci.org/vipally/binary) [![Coverage Status](https://coveralls.io/repos/github/vipally/binary/badge.svg?branch=master)](https://coveralls.io/github/vipally/binary?branch=master) [![GoDoc](https://godoc.org/github.com/vipally/binary?status.svg)](https://godoc.org/github.com/vipally/binary) ![Version](https://img.shields.io/badge/version-1.0.0.final-green.svg)
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
## 	Use RegStruct to improving struct encoding/decoding efficiency.
	type someRegedStruct struct {
		A int `binary:"ignore"`
		B string
		C uint
	}
	binary.RegStruct((*someRegedStruct)(nil))

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
	
# 9. Test reults.
## Enncoding size(see example of Sizeof).
	Encoding bytes is much shorter than std.binary and gob.

	var s struct {
		Int8        int8
		Int16       int16
		Int32       int32
		Int64       int64
		Uint8       uint8
		Uint16      uint16
		Uint32      uint32
		Uint64      uint64
		Float32     float32
		Float64     float64
		Complex64   complex64
		Complex128  complex128
		Array       [10]uint8
		Bool        bool
		BoolArray   [100]bool
		Uint32Array [10]uint32
	}

	Output:
	Sizeof(s)  = 133
	std.Size(s)= 217
	gob.Size(s)= 415
	
## Benchmark test result.
	Pack/Unpack unregisted struct is a little slower(80%) than std.binary.Read/std.binary.Write
	Pack/Unpack registed struct is much faster(200%) than std.binary.Read/std.binary.Write
	Pack/Unpack int array is much faster(230%) than std.binary.Read/std.binary.Write
	Pack/Unpack string is a little faster(140%) than gob.Encode/gob.Decode

	BenchmarkGobEncodeStruct     	 1000000	      1172 ns/op	 300.32 MB/s
	BenchmarkStdWriteStruct      	 1000000	      3154 ns/op	  23.78 MB/s
	BenchmarkWriteStruct         	  500000	      3954 ns/op	  18.71 MB/s
	BenchmarkWriteRegedStruct    	 1000000	      1627 ns/op	  45.48 MB/s
	BenchmarkPackStruct          	  500000	      3642 ns/op	  20.32 MB/s
	BenchmarkPackRegedStruct     	 1000000	      1542 ns/op	  47.99 MB/s
	BenchmarkGobDecodeStruct     	 3000000	       651 ns/op	 540.40 MB/s
	BenchmarkStdReadStruct       	 1000000	      2008 ns/op	  37.35 MB/s
	BenchmarkReadStruct          	  500000	      2386 ns/op	  31.01 MB/s
	BenchmarkReadRegedStruct     	 1000000	      1194 ns/op	  61.97 MB/s
	BenchmarkUnackStruct         	 1000000	      2293 ns/op	  32.27 MB/s
	BenchmarkUnpackRegedStruct   	 2000000	       935 ns/op	  79.14 MB/s
	BenchmarkGobEncodeInt1000    	  100000	     22871 ns/op	 219.58 MB/s
	BenchmarkStdWriteInt1000     	   30000	     49502 ns/op	  80.80 MB/s
	BenchmarkWriteInt1000        	   50000	     26001 ns/op	 153.91 MB/s
	BenchmarkPackInt1000         	  100000	     21601 ns/op	 185.27 MB/s
	BenchmarkStdReadInt1000      	   30000	     44035 ns/op	  90.84 MB/s
	BenchmarkReadInt1000         	   50000	     29761 ns/op	 134.47 MB/s
	BenchmarkUnackInt1000        	   50000	     25001 ns/op	 160.07 MB/s
	BenchmarkGobEncodeString     	 3000000	       444 ns/op	 301.56 MB/s
	BenchmarkStdWriteString 		unsupported 
	BenchmarkWriteString         	 3000000	       455 ns/op	 285.49 MB/s
	BenchmarkPackString          	 5000000	       337 ns/op	 385.05 MB/s
	BenchmarkGobDecodeString     	 1000000	      1266 ns/op	 105.84 MB/s
	BenchmarkStdReadString 			unsupported 
	BenchmarkReadString          	 3000000	       550 ns/op	 236.06 MB/s
	BenchmarkUnackString         	 5000000	       298 ns/op	 435.92 MB/s
	
# With full-datatype support like gob.
	Followed struct fullStruct is an aviable type that contains all supported types.
	Which is not valid for std.binary.
	
	type fullStruct struct {
		Int8       int8
		Int16      int16
		Int32      int32
		Int64      int64
		Uint8      uint8
		Uint16     uint16
		Uint32     uint32
		Uint64     uint64
		Float32    float32
		Float64    float64
		Complex64  complex64
		Complex128 complex128
		Array      [4]uint8
		Bool       bool
		BoolArray  [9]bool
	
		LittleStruct  littleStruct
		PLittleStruct *littleStruct
		String        string
		PString       *string
		PInt32        *int32
		Slice         []*littleStruct
		PSlice        *[]*string
		Float64Slice  []float64
		BoolSlice     []bool
		Uint32Slice   []uint32
		Map           map[string]*littleStruct
		Map2          map[string]uint16
		IntSlice      []int
		UintSlice     []uint
	}

# Do not supported types.
	Followed struct TDoNotSupport is an invalid type that every field is invalid.

	type TDoNotSupport struct {
		Uintptr       uintptr
		UnsafePointer unsafe.Pointer
		Ch            chan bool
		Map           map[uintptr]uintptr
		Map2          map[int]uintptr
		Map3          map[uintptr]int
		Slice         []uintptr
		Array         [2]uintptr
		Array2        [2][2]uintptr
		Array3        [2]struct{ A uintptr }
		Func          func()
		Struct        struct {
			PStruct *struct {
				PPUintptr **uintptr
			}
		}
	}

