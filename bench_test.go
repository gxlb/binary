package binary

import (
	"bytes"
	std "encoding/binary"
	"encoding/gob"
	"os"
	"reflect"
	"runtime/pprof"
	"testing"
)

const (
	prof = false //run pprof
)

type regedStruct struct {
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
	BoolArray  [4]bool
}

var (
	buff   = make([]byte, 8192)
	buffer = bytes.NewBuffer(buff[:0])

	wStruct Struct

	u32Array1000  [1000]uint32
	u32Array1000W [1000]uint32

	caseStdReadWrite string

	str  = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	strW string
)

func init() {
	RegisterType((*regedStruct)(nil))
	for i := len(u32Array1000) - 1; i >= 0; i-- {
		u32Array1000[i] = uint32(i)*7368787 + 2750159 //rand number
	}

	if prof {
		f, err := os.Create("b.prof")
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
	}
}

//////////////////////////////////////////////////////////////////Struct
//func BenchmarkGobEncodeStruct(b *testing.B) {
//	data := _struct
//	testBenchGobEncode(b, data, "BenchmarkGobEncodeStruct")
//}
func BenchmarkStdWriteStruct(b *testing.B) {
	data := _struct
	testBenchStdWrite(b, data, "BenchmarkStdWriteStruct")
}
func BenchmarkWriteStruct(b *testing.B) {
	data := _struct
	testBenchWrite(b, data, "BenchmarkWriteStruct")
}
func BenchmarkWriteRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	testBenchWrite(b, data, "BenchmarkWriteRegedStruct")
}
func BenchmarkEncodeStruct(b *testing.B) {
	data := _struct
	testBenchEncode(b, data, "BenchmarkEncodeStruct")
}
func BenchmarkEncodeRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	testBenchEncode(b, data, "BenchmarkEncodeRegedStruct")
}

//func BenchmarkGobDecodeStruct(b *testing.B) {
//	data := _struct
//	testBenchGobDecode(b, &data, &wStruct, "BenchmarkGobDecodeStruct")
//}
func BenchmarkStdReadStruct(b *testing.B) {
	data := _struct
	testBenchStdRead(b, &data, &wStruct, "BenchmarkStdReadStruct")
}

func BenchmarkReadStruct(b *testing.B) {
	data := _struct
	testBenchRead(b, &data, &wStruct, "BenchmarkReadStruct")
}
func BenchmarkReadRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	dataC := regedStruct(wStruct)
	testBenchRead(b, &data, &dataC, "BenchmarkReadRegedStruct")
}
func BenchmarkDecodeStruct(b *testing.B) {
	data := _struct
	testBenchDecode(b, &data, &wStruct, "BenchmarkDecodeStruct")
}
func BenchmarkDecodeRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	dataC := regedStruct(wStruct)
	testBenchDecode(b, &data, &dataC, "BenchmarkDecodeRegedStruct")
}

//////////////////////////////////////////////////////////////////Int1000
//func BenchmarkGobEncodeInt1000(b *testing.B) {
//	data := u32Array1000
//	testBenchGobEncode(b, data, "BenchmarkGobEncodeInt1000")
//}
func BenchmarkStdWriteInt1000(b *testing.B) {
	data := u32Array1000
	testBenchStdWrite(b, data, "BenchmarkStdWriteInt1000")
}
func BenchmarkWriteInt1000(b *testing.B) {
	data := u32Array1000
	testBenchWrite(b, data, "BenchmarkWriteInt1000")
}
func BenchmarkEncodeInt1000(b *testing.B) {
	data := u32Array1000
	testBenchEncode(b, data, "BenchmarkEncodeInt1000")
} //BUG: this case will fail
//func BenchmarkGobDecodeInt1000(b *testing.B) {
//	data := u32Array1000
//	testBenchGobDecode(b, &data, &u32Array1000W, "BenchmarkGobDecodeInt1000")
//}
func BenchmarkStdReadInt1000(b *testing.B) {
	data := u32Array1000
	testBenchStdRead(b, &data, &u32Array1000W, "BenchmarkStdReadInt1000")
}
func BenchmarkReadInt1000(b *testing.B) {
	data := u32Array1000
	testBenchRead(b, &data, &u32Array1000W, "BenchmarkReadInt1000")
}
func BenchmarkDecodeInt1000(b *testing.B) {
	data := u32Array1000
	testBenchDecode(b, &data, &u32Array1000W, "BenchmarkUnackInt1000")
}

//////////////////////////////////////////////////////////////////String
//func BenchmarkGobEncodeString(b *testing.B) {
//	data := str
//	testBenchGobEncode(b, data, "BenchmarkGobEncodeString")
//}
func BenchmarkStdWriteString(b *testing.B) {
	data := str
	testBenchStdWrite(b, data, "BenchmarkStdWriteString")
}
func BenchmarkWriteString(b *testing.B) {
	data := str
	testBenchWrite(b, data, "BenchmarkWriteString")
}
func BenchmarkEncodeString(b *testing.B) {
	data := str
	testBenchEncode(b, data, "BenchmarkEncodeString")
}

//func BenchmarkGobDecodeString(b *testing.B) {
//	data := str
//	testBenchGobDecode(b, &data, &strW, "BenchmarkGobDecodeString")
//}
func BenchmarkStdReadString(b *testing.B) {
	data := str
	testBenchStdRead(b, &data, &strW, "BenchmarkStdReadString")
}
func BenchmarkReadString(b *testing.B) {
	data := str
	testBenchRead(b, &data, &strW, "BenchmarkReadString")
}
func BenchmarkDecodeString(b *testing.B) {
	data := str
	testBenchDecode(b, &data, &strW, "BenchmarkUnackString")
}

//func newSame(v reflect.Value) (value reflect.Value) {
//	vv := reflect.Indirect(v)
//	t := vv.Type()
//	switch t.Kind() {
//	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16,
//		reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64,
//		reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64,
//		reflect.Complex128, reflect.String, reflect.Array, reflect.Struct:
//		value = reflect.New(t)
//	case reflect.Slice:
//		value = reflect.MakeSlice(t, 0, 0).Addr() //make a default slice
//	}
//	return
//}

func testBenchGobEncode(b *testing.B, data interface{}, caseName string) {
	buffer.Reset()
	coder := gob.NewEncoder(buffer)
	err := coder.Encode(data)
	b.SetBytes(int64(buffer.Len()))
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		coder.Encode(data)
	}
	b.StopTimer()
}
func testBenchStdWrite(b *testing.B, data interface{}, caseName string) {
	s := std.Size(data)
	if s <= 0 {
		if caseStdReadWrite != caseName {
			caseStdReadWrite = caseName
			println(caseName, "unsupported ")
		}
		return
	}
	buffer.Reset()
	b.SetBytes(int64(s))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		std.Write(buffer, std.LittleEndian, data)
	}
	b.StopTimer()
}
func testBenchWrite(b *testing.B, data interface{}, caseName string) {
	b.SetBytes(int64(SizeX(data, false)))
	buffer.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		Write(buffer, LittleEndian, data)
	}
	b.StopTimer()
}
func testBenchEncode(b *testing.B, data interface{}, caseName string) {
	b.SetBytes(int64(SizeX(data, false)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		EncodeX(data, buff, false)
	}
	b.StopTimer()
}

func testBenchGobDecode(b *testing.B, data, w interface{}, caseName string) {
	bsr := &byteSliceReader{}
	buffer.Reset()
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(buffer.Len()))

	b.ResetTimer()
	buf := buffer.Bytes()
	bsr.remain = buf
	decoder := gob.NewDecoder(bsr)
	decoder.Decode(w)
	for i := 0; i < b.N; i++ {
		bsr.remain = buf
		decoder.Decode(w)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, w) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, w, data)
	}
}
func testBenchStdRead(b *testing.B, data, w interface{}, caseName string) {
	s := std.Size(data)
	if s <= 0 {
		if caseStdReadWrite != caseName {
			caseStdReadWrite = caseName
			println(caseName, "unsupported ")
		}
		return
	}
	bsr := &byteSliceReader{}
	buffer.Reset()
	err := std.Write(buffer, std.LittleEndian, data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buffer.Bytes())))
	b.ResetTimer()
	bsr.remain = buffer.Bytes()
	std.Read(bsr, std.LittleEndian, w)
	for i := 0; i < b.N; i++ {
		bsr.remain = buffer.Bytes()
		std.Read(bsr, std.LittleEndian, w)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, w) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, w, data)
	}
}
func testBenchRead(b *testing.B, data, w interface{}, caseName string) {
	bsr := &byteSliceReader{}
	buffer.Reset()
	err := Write(buffer, DefaultEndian, data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buffer.Bytes())))

	b.ResetTimer()
	bsr.remain = buffer.Bytes()
	Read(bsr, DefaultEndian, w)
	for i := 0; i < b.N; i++ {
		bsr.remain = buffer.Bytes()
		Read(bsr, DefaultEndian, w)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, w) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, w, data)
	}
}
func testBenchDecode(b *testing.B, data, w interface{}, caseName string) {
	buf, err := EncodeX(data, buff, false)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DecodeX(buf, w, false)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, w) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, w, data)
	}
}

func TestBenchmarkEnd(t *testing.T) {
	if prof {
		pprof.StopCPUProfile()
	}
}
