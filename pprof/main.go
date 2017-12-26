// usage:
//pprof.exe -n=x 1~4
//go tool pprof gb.prof
//png
//quit

package main

func main() {}

//import (
//	"bytes"
//	std "encoding/binary"
//	"flag"
//	"fmt"
//	"os"
//	"runtime/pprof"
//	"time"

//	. "github.com/vipally/binary"
//)

//var (
//	u32Array1000     [1000]uint32
//	u32Array1000W    [1000]uint32
//	buff             = make([]byte, 8192)
//	buffer           = bytes.NewBuffer(buff[:0])
//	N                = 1000000
//	caseStdReadWrite string
//	str              = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
//	strW             string
//)

//func init() {
//	for i := len(u32Array1000) - 1; i >= 0; i-- {
//		u32Array1000[i] = uint32(i)*7368787 + 2750159 //rand number
//	}
//}

//func main() {
//	f, err := os.Create("gb.prof")
//	if err != nil {
//		panic(err)
//	}
//	pprof.StartCPUProfile(f)
//	defer pprof.StopCPUProfile()
//	n := flag.Int("n", 0, "sel number 1~4")
//	flag.Parse()

//	start := time.Now()
//	switch *n {
//	case 1:
//		testBenchEncode(u32Array1000, "BenchmarkEncodeInt1000")
//	case 2:
//		testBenchDecode(&u32Array1000, &u32Array1000W, "BenchmarkUnackInt1000")
//	case 3:
//		testBenchStdWrite(u32Array1000, "BenchmarkStdWriteInt1000")
//	case 4:
//		testBenchRead(&u32Array1000, &u32Array1000W, "BenchmarkReadInt1000")
//	case 5:
//		testBenchEncode(str, "BenchmarkEncodeString")
//	case 6:
//		testBenchDecode(&str, &strW, "BenchmarkUnackString")
//	case 7:
//		testBenchStdWrite(str, "BenchmarkStdWriteString")
//	case 8:
//		testBenchStdRead(&str, &strW, "BenchmarkStdReadString")
//	}
//	dur := time.Now().Sub(start)
//	fmt.Println("finish dur=", dur)
//}

////func testBenchGobEncode(data interface{}, caseName string) {
////	buffer.Reset()
////	coder := gob.NewEncoder(buffer)
////	err := coder.Encode(data)
////	b.SetBytes(int64(buffer.Len()))
////	if err != nil {
////		b.Error(err)
////	}
////	b.ResetTimer()
////	for i := 0; i < b.N; i++ {
////		buffer.Reset()
////		coder.Encode(data)
////	}
////	b.StopTimer()
////}
//func testBenchStdWrite(data interface{}, caseName string) {
//	s := std.Size(data)
//	if s <= 0 {
//		if caseStdReadWrite != caseName {
//			caseStdReadWrite = caseName
//			println(caseName, "unsupported ")
//		}
//		return
//	}
//	buffer.Reset()

//	for i := 0; i < N; i++ {
//		buffer.Reset()
//		std.Write(buffer, std.LittleEndian, data)
//	}
//}
//func testBenchWrite(data interface{}, caseName string) {
//	buffer.Reset()
//	for i := 0; i < N; i++ {
//		buffer.Reset()
//		Write(buffer, LittleEndian, data)
//	}
//}
//func testBenchEncode(data interface{}, caseName string) {
//	for i := 0; i < N; i++ {
//		EncodeX(data, buff, false)
//	}
//}

////func testBenchGobDecode(data, w interface{}, caseName string) {
////	bsr := &byteSliceReader{}
////	buffer.Reset()
////	encoder := gob.NewEncoder(buffer)
////	err := encoder.Encode(data)
////	if err != nil {
////		b.Error(caseName, err)
////	}
////	b.SetBytes(int64(buffer.Len()))

////	b.ResetTimer()
////	buf := buffer.Bytes()
////	bsr.remain = buf
////	decoder := gob.NewDecoder(bsr)
////	decoder.Decode(w)
////	for i := 0; i < b.N; i++ {
////		bsr.remain = buf
////		decoder.Decode(w)
////	}
////	b.StopTimer()
////	if b.N > 0 && !reflect.DeepEqual(data, w) {
////		b.Fatalf("%s doesn't match:\ngot  %#v;\nwant %#v", caseName, w, data)
////	}
////}
//func testBenchStdRead(data, w interface{}, caseName string) {
//	s := std.Size(data)
//	if s <= 0 {
//		if caseStdReadWrite != caseName {
//			caseStdReadWrite = caseName
//			println(caseName, "unsupported ")
//		}
//		return
//	}
//	bsr := &byteSliceReader{}
//	buffer.Reset()
//	err := std.Write(buffer, std.LittleEndian, data)
//	if err != nil {
//		fmt.Println(caseName, err)
//	}

//	bsr.remain = buffer.Bytes()
//	std.Read(bsr, std.LittleEndian, w)
//	for i := 0; i < N; i++ {
//		bsr.remain = buffer.Bytes()
//		std.Read(bsr, std.LittleEndian, w)
//	}

//}
//func testBenchRead(data, w interface{}, caseName string) {
//	bsr := &byteSliceReader{}
//	buffer.Reset()
//	err := Write(buffer, DefaultEndian, data)
//	if err != nil {
//		fmt.Println(caseName, err)
//	}

//	bsr.remain = buffer.Bytes()
//	Read(bsr, DefaultEndian, w)
//	for i := 0; i < N; i++ {
//		bsr.remain = buffer.Bytes()
//		Read(bsr, DefaultEndian, w)
//	}
//}
//func testBenchDecode(data, w interface{}, caseName string) {
//	buf, err := EncodeX(data, buff, false)
//	if err != nil {
//		fmt.Println(caseName, err)
//	}

//	for i := 0; i < N; i++ {
//		DecodeX(buf, w, false)
//	}
//}

//type byteSliceReader struct {
//	remain []byte
//}

//func (br *byteSliceReader) Read(p []byte) (int, error) {
//	n := copy(p, br.remain)
//	br.remain = br.remain[n:]
//	return n, nil
//}
