package binary

///*
import (
	"bytes"
	std "encoding/binary"
	"encoding/gob"
	"reflect"
	"testing"
)

type regedStruct Struct

var (
	buff   = make([]byte, 8192)
	buffer = bytes.NewBuffer(buff[:0])

	u32Array1000 [1000]uint32

//	str               = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
//	cpx128 complex128 = complex(111.5, 555.5)
)

func init() {
	RegistStruct((*regedStruct)(nil))
	for i := len(u32Array1000) - 1; i >= 0; i-- {
		u32Array1000[i] = uint32(i)*7368787 + 2750159 //rand number
	}
}

func BenchmarkGobEncodeStruct(b *testing.B) {
	data := _struct
	testBenchGobEncode(b, data, "BenchmarkGobEncodeStruct")
}
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
func BenchmarkPackStruct(b *testing.B) {
	data := _struct
	testBenchPack(b, data, "BenchmarkPackStruct")
}
func BenchmarkPackRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	testBenchPack(b, data, "BenchmarkPackRegedStruct")
}
func BenchmarkGobDecodeStruct(b *testing.B) {
	data := _struct
	testBenchGobDecode(b, &data, "BenchmarkGobDecodeStruct")
}
func BenchmarkStdReadStruct(b *testing.B) {
	data := _struct
	testBenchStdRead(b, &data, "BenchmarkStdReadStruct")
}
func BenchmarkReadStruct(b *testing.B) {
	data := _struct
	testBenchRead(b, &data, "BenchmarkReadStruct")
}
func BenchmarkReadRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	testBenchRead(b, &data, "BenchmarkReadRegedStruct")
}
func BenchmarkUnackStruct(b *testing.B) {
	data := _struct
	testBenchUnpack(b, &data, "BenchmarkUnackStruct")
}
func BenchmarkUnpackRegedStruct(b *testing.B) {
	data := regedStruct(_struct)
	testBenchUnpack(b, &data, "BenchmarkUnpackRegedStruct")
}

/////////////////////////////////
func BenchmarkGobEncodeInt1000(b *testing.B) {
	data := u32Array1000
	testBenchGobEncode(b, data, "BenchmarkGobEncodeInt1000")
}
func BenchmarkStdWriteInt1000(b *testing.B) {
	data := u32Array1000
	testBenchStdWrite(b, data, "BenchmarkStdWriteInt1000")
}
func BenchmarkWriteInt1000(b *testing.B) {
	data := u32Array1000
	testBenchWrite(b, data, "BenchmarkWriteInt1000")
}
func BenchmarkPackInt1000(b *testing.B) {
	data := u32Array1000
	testBenchPack(b, data, "BenchmarkPackInt1000")
}

//func BenchmarkGobDecodeInt1000(b *testing.B) {
//	data := u32Array1000
//	testBenchGobDecode(b, &data, "BenchmarkGobDecodeInt1000")
//}
func BenchmarkStdReadInt1000(b *testing.B) {
	data := u32Array1000
	testBenchStdRead(b, &data, "BenchmarkStdReadInt1000")
}
func BenchmarkReadInt1000(b *testing.B) {
	data := u32Array1000
	testBenchRead(b, &data, "BenchmarkReadInt1000")
}
func BenchmarkUnackInt1000(b *testing.B) {
	data := u32Array1000
	testBenchUnpack(b, &data, "BenchmarkUnackInt1000")
}

////////////////
//func BenchmarkReadStruct1(b *testing.B) {
//	testBenchRead(b, &_struct, "struct")
//}

//func BenchmarkUnpackStruct1(b *testing.B) {
//	buff, _ := Pack(_struct, nil)
//	b.SetBytes(int64(sizeofValue(reflect.ValueOf(_struct))))
//	var t Struct
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		Unpack(buff, &t)
//	}
//	b.StopTimer()
//	if b.N > 0 && !reflect.DeepEqual(_struct, t) {
//		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
//	}
//}

//func BenchmarkReadString(b *testing.B) {
//	bsr := &byteSliceReader{}
//	var buf bytes.Buffer
//	Write(&buf, BigEndian, str)
//	b.SetBytes(int64(buf.Len()))
//	var t string
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		bsr.remain = buf.Bytes()
//		Read(bsr, BigEndian, &t)
//	}
//	b.StopTimer()
//	if b.N > 0 && !reflect.DeepEqual(str, t) {
//		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
//	}
//}
//func BenchmarkUnpackString(b *testing.B) {
//	buff, _ := Pack(str, nil)
//	b.SetBytes(int64(len(buff)))
//	var t string
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		Unpack(buff, &t)
//	}
//	b.StopTimer()
//	if b.N > 0 && !reflect.DeepEqual(str, t) {
//		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
//	}
//}

//func testBenchReadStd(b *testing.B, data interface{}, caseName string) {
//	bsr := &byteSliceReader{}
//	var buf bytes.Buffer
//	Write(&buf, BigEndian, data)
//	b.SetBytes(int64(sizeofValue(reflect.ValueOf(_struct))))
//	t := reflect.New(reflect.Indirect(reflect.ValueOf(data)).Type())
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		bsr.remain = buf.Bytes()
//		Read(bsr, BigEndian, &t)
//	}
//	b.StopTimer()
//	if b.N > 0 && !reflect.DeepEqual(_struct, t) {
//		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
//	}
//}

func newSame(v reflect.Value) (value reflect.Value) {
	vv := reflect.Indirect(v)
	t := vv.Type()
	switch t.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8, reflect.Int16,
		reflect.Uint16, reflect.Int32, reflect.Uint32, reflect.Int64,
		reflect.Uint64, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.String, reflect.Array, reflect.Struct:
		value = reflect.New(t)
	case reflect.Slice:
		value = reflect.MakeSlice(t, 0, 0).Addr() //make a default slice
	}
	return
}

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
	buffer.Reset()
	b.SetBytes(int64(std.Size(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		std.Write(buffer, std.LittleEndian, data)
	}
	b.StopTimer()
}
func testBenchWrite(b *testing.B, data interface{}, caseName string) {
	b.SetBytes(int64(Sizeof(data)))
	buffer.Reset()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		Write(buffer, std.LittleEndian, data)
	}
	b.StopTimer()
}
func testBenchPack(b *testing.B, data interface{}, caseName string) {
	b.SetBytes(int64(Sizeof(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(data, buff)
	}
	b.StopTimer()
}

func testBenchGobDecode(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	buffer.Reset()
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(buffer.Len()))
	newValue := newSame(reflect.ValueOf(data))
	t := newValue.Interface()
	b.ResetTimer()
	bsr.remain = buffer.Bytes()
	decoder := gob.NewDecoder(bsr)
	decoder.Decode(t)
	for i := 0; i < b.N; i++ {
		bsr.remain = buffer.Bytes()
		decoder.Decode(t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, t) {
		//b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, t, data)
	}
}
func testBenchStdRead(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	buffer.Reset()
	err := std.Write(buffer, std.LittleEndian, data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buffer.Bytes())))
	newValue := newSame(reflect.ValueOf(data))
	t := newValue.Interface()
	b.ResetTimer()
	bsr.remain = buffer.Bytes()
	std.Read(bsr, std.LittleEndian, t)
	for i := 0; i < b.N; i++ {
		bsr.remain = buffer.Bytes()
		std.Read(bsr, std.LittleEndian, t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, t) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, t, data)
	}
}
func testBenchRead(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	buffer.Reset()
	err := Write(buffer, DefaultEndian, data)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buffer.Bytes())))
	newValue := newSame(reflect.ValueOf(data))
	t := newValue.Interface()
	b.ResetTimer()
	bsr.remain = buffer.Bytes()
	Read(bsr, DefaultEndian, t)
	for i := 0; i < b.N; i++ {
		bsr.remain = buffer.Bytes()
		Read(bsr, DefaultEndian, t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, t) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, t, data)
	}
}
func testBenchUnpack(b *testing.B, data interface{}, caseName string) {
	buf, err := Pack(data, buff)
	if err != nil {
		b.Error(caseName, err)
	}
	b.SetBytes(int64(len(buf)))
	newValue := newSame(reflect.ValueOf(data))
	t := newValue.Interface()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unpack(buf, t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, t) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, t, data)
	}
}

//*/
