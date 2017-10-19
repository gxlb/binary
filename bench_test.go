package binary

///*
import (
	"bytes"
	std "encoding/binary"
	"encoding/gob"
	"reflect"
	"testing"
)

type regStruct Struct

func init() {
	RegistStruct((*regStruct)(nil))
}

func BenchmarkGobEncode(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	coder := gob.NewEncoder(buffer)
	err := coder.Encode(_struct)
	b.SetBytes(int64(buffer.Len()))
	if err != nil {
		b.Error(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		coder.Encode(_struct)
	}
	b.StopTimer()
	//println(len(buffer.Bytes()))
}
func BenchmarkStdWriteStruct(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	b.SetBytes(int64(std.Size(_struct)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		std.Write(buffer, std.LittleEndian, _struct)
	}
	b.StopTimer()
}
func BenchmarkWriteStruct(b *testing.B) {
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	b.SetBytes(int64(Sizeof(_struct)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		Write(buffer, std.LittleEndian, _struct)
	}
	b.StopTimer()
}
func BenchmarkWriteRegedStruct(b *testing.B) {
	data := regStruct(_struct)
	buffer := bytes.NewBuffer(make([]byte, 0, 1024))
	b.SetBytes(int64(Sizeof(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		Write(buffer, std.LittleEndian, data)
	}
	b.StopTimer()
}

func BenchmarkPackStruct(b *testing.B) {
	buffer := make([]byte, 0, 2048)
	b.SetBytes(int64(Sizeof(_struct)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(_struct, buffer)
	}
	b.StopTimer()
	//fmt.Println(err, len(buffer.Bytes()))
}

func BenchmarkPackRegedStruct(b *testing.B) {
	data := regStruct(_struct)
	buffer := make([]byte, 0, 2048)
	b.SetBytes(int64(Sizeof(data)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(data, buffer)
	}
	b.StopTimer()
	//fmt.Println(err, len(buffer.Bytes()))
}

func BenchmarkReadStruct1(b *testing.B) {
	testBenchRead(b, &_struct, "struct")
}

func BenchmarkUnpackStruct1(b *testing.B) {
	buff, _ := Pack(_struct, nil)
	b.SetBytes(int64(sizeofValue(reflect.ValueOf(_struct))))
	var t Struct
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unpack(buff, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(_struct, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
	}
}

var (
	str               = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	cpx128 complex128 = complex(111.5, 555.5)
)

func BenchmarkReadString(b *testing.B) {
	bsr := &byteSliceReader{}
	var buf bytes.Buffer
	Write(&buf, BigEndian, str)
	b.SetBytes(int64(buf.Len()))
	var t string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bsr.remain = buf.Bytes()
		Read(bsr, BigEndian, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(str, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
	}
}
func BenchmarkUnpackString(b *testing.B) {
	buff, _ := Pack(str, nil)
	b.SetBytes(int64(len(buff)))
	var t string
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unpack(buff, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(str, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
	}
}

func testBenchReadStd(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	var buf bytes.Buffer
	Write(&buf, BigEndian, data)
	b.SetBytes(int64(sizeofValue(reflect.ValueOf(_struct))))
	t := reflect.New(reflect.Indirect(reflect.ValueOf(data)).Type())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bsr.remain = buf.Bytes()
		Read(bsr, BigEndian, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(_struct, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, _struct)
	}
}

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

func testBenchRead(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	var buf bytes.Buffer
	Write(&buf, BigEndian, data)
	b.SetBytes(int64(Size(data)))
	newValue := newSame(reflect.ValueOf(data))
	t := newValue.Interface()
	//tt := newValue.Elem().Interface()
	b.ResetTimer()
	bsr.remain = buf.Bytes()
	Read(bsr, BigEndian, t)
	//	fmt.Printf("%#v\n%+v\n", data, buf.Bytes())
	//	fmt.Printf("%#v\n", t)
	//	fmt.Printf("%#v\n", tt)
	for i := 0; i < b.N; i++ {
		bsr.remain = buf.Bytes()
		Read(bsr, BigEndian, t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(data, t) {
		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, t, data)
	}
}
func testBenchUnack(b *testing.B, data interface{}, caseName string)    {}
func testBenchWriteStd(b *testing.B, data interface{}, caseName string) {}
func testBenchWrite(b *testing.B, data interface{}, caseName string)    {}
func testBenchPack(b *testing.B, data interface{}, caseName string)     {}

func BenchmarkEncoder(b *testing.B) {

}

//*/
