package binary

/*
import (
	"bytes"
	//"fmt"
	"reflect"
	"testing"
)

func BenchmarkReadStruct1(b *testing.B) {
	testBenchRead(b, &s, "struct")
}

func BenchmarkUnpackStruct1(b *testing.B) {
	buff, _ := Pack(s, nil)
	b.SetBytes(int64(sizeofValue(reflect.ValueOf(s))))
	var t Struct
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unpack(buff, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(s, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, s)
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
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, s)
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
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, s)
	}
}

func testBenchReadStd(b *testing.B, data interface{}, caseName string) {
	bsr := &byteSliceReader{}
	var buf bytes.Buffer
	Write(&buf, BigEndian, data)
	b.SetBytes(int64(sizeofValue(reflect.ValueOf(s))))
	t := reflect.New(reflect.Indirect(reflect.ValueOf(data)).Type())
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bsr.remain = buf.Bytes()
		Read(bsr, BigEndian, &t)
	}
	b.StopTimer()
	if b.N > 0 && !reflect.DeepEqual(s, t) {
		b.Fatalf("struct doesn't match:\ngot  %v;\nwant %v", t, s)
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

//*/
