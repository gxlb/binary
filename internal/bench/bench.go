//bench test case for binary and std

package bench

import (
	"bufio"
	"bytes"
	std "encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/vipally/binary"
	"github.com/vipally/binary/random"
)

var (
	full FullStruct
)

type benchType byte

const (
	BenchStdWrite benchType = iota
	BenchStdRead
	BenchEncode
	BenchDecode
	BenchGobEncode
	BenchGobDecode

	benchDoCnt = 100000
)

func (bench benchType) String() string {
	switch bench {
	case BenchStdWrite:
		return "BenchStdWrite"
	case BenchStdRead:
		return "BenchStdRead"
	case BenchEncode:
		return "BenchEncode"
	case BenchDecode:
		return "BenchDecode"
	case BenchGobEncode:
		return "BenchGobEncode"
	case BenchGobDecode:
		return "BenchGobDecode"
	}
	panic("undefined benchType")
}

var (
	buff   = make([]byte, 8192)
	buffer = bytes.NewBuffer(buff[:0])
)

func BenchCases() []*BenchCase {
	return cases
}

type Speed float64

func (s Speed) String() string {
	if s < 0 {
		return "-"
	}
	if s >= 1000000 {
		return fmt.Sprintf("%.2fMOp/s", s/1000000)
	}

	if s >= 1000 {
		return fmt.Sprintf("%.2fKOp/s", s/1000)
	}
	return fmt.Sprintf("%.0fOp/s", s)
}

type Duration time.Duration

func (dur Duration) String() string {
	if dur < 0 {
		return "-"
	}
	if true || dur >= Duration(time.Second) {
		return fmt.Sprintf("%.2fs", float64(dur)/float64(time.Second))
	}
	return fmt.Sprintf("%.2fms", float64(dur)/float64(time.Millisecond))
}

type Size int

const (
	kBits = 10
	kb    = 1 << kBits
)

var bytesNames = []string{"B", "KB", "MB", "GB", "TB"} //max16EB ignore the next: "ZB", "YB", "BB"

func (s Size) String() string {
	if s <= 0 {
		return "-"
	}

	m64 := uint64(s)
	i, b := 0, uint64(0)
	for ; i < len(bytesNames); i, b = i+1, b+kBits {
		if (m64 >> b) < kb {
			break
		}
	}
	m := 1 << b
	d := float64(s) / float64(m)
	return fmt.Sprintf("%.2f%s", d, bytesNames[i])
}

var (
	UvarintCases = []uint64{
		0x0000000000000001, 0x0000000000000003, 0x0000000000000007, 0x000000000000000F,
		0x000000000000001F, 0x000000000000003F, 0x000000000000007F, 0x00000000000000FF,
		0x00000000000001FF, 0x00000000000003FF, 0x00000000000007FF, 0x0000000000000FFF,
		0x0000000000001FFF, 0x0000000000003FFF, 0x0000000000007FFF, 0x000000000000FFFF,
		0x000000000001FFFF, 0x000000000003FFFF, 0x000000000007FFFF, 0x00000000000FFFFF,
		0x00000000001FFFFF, 0x00000000003FFFFF, 0x00000000007FFFFF, 0x0000000000FFFFFF,
		0x0000000001FFFFFF, 0x0000000003FFFFFF, 0x0000000007FFFFFF, 0x000000000FFFFFFF,
		0x000000001FFFFFFF, 0x000000003FFFFFFF, 0x000000007FFFFFFF, 0x00000000FFFFFFFF,
		0x00000001FFFFFFFF, 0x00000003FFFFFFFF, 0x00000007FFFFFFFF, 0x0000000FFFFFFFFF,
		0x0000001FFFFFFFFF, 0x0000003FFFFFFFFF, 0x0000007FFFFFFFFF, 0x000000FFFFFFFFFF,
		0x000001FFFFFFFFFF, 0x000003FFFFFFFFFF, 0x000007FFFFFFFFFF, 0x00000FFFFFFFFFFF,
		0x00001FFFFFFFFFFF, 0x00003FFFFFFFFFFF, 0x00007FFFFFFFFFFF, 0x0000FFFFFFFFFFFF,
		0x0001FFFFFFFFFFFF, 0x0003FFFFFFFFFFFF, 0x0007FFFFFFFFFFFF, 0x000FFFFFFFFFFFFF,
		0x001FFFFFFFFFFFFF, 0x003FFFFFFFFFFFFF, 0x007FFFFFFFFFFFFF, 0x00FFFFFFFFFFFFFF,
		0x01FFFFFFFFFFFFFF, 0x03FFFFFFFFFFFFFF, 0x07FFFFFFFFFFFFFF, 0x0FFFFFFFFFFFFFFF,
		0x1FFFFFFFFFFFFFFF, 0x3FFFFFFFFFFFFFFF, 0x7FFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF,
	}
	UvarintStdBytes = []byte{ //cnt= 325
		0x01, 0x03, 0x07, 0x0f, 0x1f, 0x3f, 0x7f, 0xff, 0x01, 0xff, 0x03, 0xff, 0x07, 0xff, 0x0f, 0xff,
		0x1f, 0xff, 0x3f, 0xff, 0x7f, 0xff, 0xff, 0x01, 0xff, 0xff, 0x03, 0xff, 0xff, 0x07, 0xff, 0xff,
		0x0f, 0xff, 0xff, 0x1f, 0xff, 0xff, 0x3f, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff,
		0xff, 0x03, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff,
		0xff, 0x3f, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0x03,
		0xff, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff,
		0xff, 0xff, 0xff, 0x3f, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x03, 0xff, 0xff, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x03, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0f, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0x03, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x01, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x03, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x07, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x0f,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0x3f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x01,
	}
	UvarintBytes = []byte{ //cnt= 316
		0x01, 0x03, 0x07, 0x0f, 0x1f, 0x3f, 0x40, 0x7f, 0x40, 0xff, 0x41, 0xff, 0x43, 0xff, 0x47, 0xff,
		0x4f, 0xff, 0x5f, 0xff, 0x7f, 0xff, 0x80, 0xff, 0x7f, 0x80, 0xff, 0xff, 0x81, 0xff, 0xff, 0x83,
		0xff, 0xff, 0x87, 0xff, 0xff, 0x8f, 0xff, 0xff, 0x90, 0xff, 0xff, 0x1f, 0x90, 0xff, 0xff, 0x3f,
		0x90, 0xff, 0xff, 0x7f, 0x90, 0xff, 0xff, 0xff, 0x91, 0xff, 0xff, 0xff, 0x93, 0xff, 0xff, 0xff,
		0x97, 0xff, 0xff, 0xff, 0x9f, 0xff, 0xff, 0xff, 0xa0, 0xff, 0xff, 0xff, 0x1f, 0xa0, 0xff, 0xff,
		0xff, 0x3f, 0xa0, 0xff, 0xff, 0xff, 0x7f, 0xa0, 0xff, 0xff, 0xff, 0xff, 0xa1, 0xff, 0xff, 0xff,
		0xff, 0xa3, 0xff, 0xff, 0xff, 0xff, 0xa7, 0xff, 0xff, 0xff, 0xff, 0xaf, 0xff, 0xff, 0xff, 0xff,
		0xb0, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xb0, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xb0, 0xff, 0xff, 0xff,
		0xff, 0x7f, 0xb0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xb1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xb3, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xb7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xbf, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xc0, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xc0, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xc0, 0xff,
		0xff, 0xff, 0xff, 0xff, 0x7f, 0xc0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xc1, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xc3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xc7, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xcf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xd0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f,
		0xd0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xd0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f,
		0xd0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xd1, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xd3, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xd7, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xdf, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xe0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0x1f, 0xe0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3f, 0xe0, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0x7f, 0xe0, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
)

type littleUvarintCase byte
type bigUvarintCase byte
type littleUvarintCaseRead byte
type bigUvarintCaseRead byte
type littleUvarintCaseReadFile byte
type bigUvarintCaseReadFile byte

const (
	LittleUvarint         littleUvarintCase         = 0
	BigUvarint            bigUvarintCase            = 0
	LittleUvarintRead     littleUvarintCaseRead     = 0
	BigUvarintRead        bigUvarintCaseRead        = 0
	LittleUvarintReadFile littleUvarintCaseReadFile = 0
	BigUvarintReadFile    bigUvarintCaseReadFile    = 0
	delta                                           = 16
	deltaBig                                        = delta * 0x100000000
)

func DoBenchUvarint(bench benchType, data interface{}, doCnt int) (t Duration, speed Speed, size Size) {
	start := time.Now()
	byteNum := 0
	var buf [10]byte
	buff := buf[0:]
	switch bench {
	case BenchStdWrite:
		switch datax := data.(type) {
		case []uint64:
			for i := 0; i < doCnt; i++ {
				byteNum = 0
				for _, x := range datax {
					byteNum += std.PutUvarint(buff, x)
				}
			}
		case littleUvarintCase:
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				byteNum += std.PutUvarint(buff, x)
				std.Uvarint(buff)
			}
		case littleUvarintCaseRead:
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				byteNum += std.PutUvarint(buff, x)
				reader := binary.BytesReader(buff)
				std.ReadUvarint(&reader)
			}
		case littleUvarintCaseReadFile:
			file := "stduvarintlittle.hex"
			f, _ := os.Create(file)
			w := bufio.NewWriter(f)
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				n := std.PutUvarint(buff, x)
				byteNum += n
				w.Write(buff[:n])
			}
			f.Sync()
			f.Close()
			f, _ = os.Open(file)
			r := bufio.NewReader(f)
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				std.ReadUvarint(r)
			}
			f.Close()
			os.Remove(file)
		case bigUvarintCase:
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				byteNum += std.PutUvarint(buff, x)
				std.Uvarint(buff)
			}
		case bigUvarintCaseRead:
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				byteNum += std.PutUvarint(buff, x)
				reader := binary.BytesReader(buff)
				std.ReadUvarint(&reader)
			}
		case bigUvarintCaseReadFile:
			file := "stduvarintbig.hex"
			f, _ := os.Create(file)
			w := bufio.NewWriter(f)
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				n := std.PutUvarint(buff, x)
				byteNum += n
				w.Write(buff[:n])
			}
			f.Sync()
			f.Close()
			f, _ = os.Open(file)
			r := bufio.NewReader(f)
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				std.ReadUvarint(r)
			}
			f.Close()
			os.Remove(file)
		}

	case BenchStdRead:
		switch datax := data.(type) {
		case []byte:
			for i := 0; i < doCnt; i++ {
				byteNum = 0
				for byteNum < len(datax) {
					_, n := std.Uvarint(datax[byteNum:])
					byteNum += n
				}
			}
		case littleUvarintCase:
		case bigUvarintCase:
		}

	case BenchEncode:
		switch datax := data.(type) {
		case []uint64:
			for i := 0; i < doCnt; i++ {
				byteNum = 0
				for _, x := range datax {
					byteNum += binary.PutUvarint(buff, x)
				}
			}
		case littleUvarintCase:
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				byteNum += binary.PutUvarint(buff, x)
				binary.Uvarint(buff)
			}
		case littleUvarintCaseRead:
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				byteNum += binary.PutUvarint(buff, x)
				reader := binary.BytesReader(buff)
				binary.ReadUvarint(&reader)
			}
		case littleUvarintCaseReadFile:
			file := "uvarintlittle.hex"
			f, _ := os.Create(file)
			w := bufio.NewWriter(f)
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				n := binary.PutUvarint(buff, x)
				byteNum += n
				w.Write(buff[:n])
			}
			f.Sync()
			f.Close()
			f, _ = os.Open(file)
			r := bufio.NewReader(f)
			for x := uint64(0); x <= 0xFFFFFFFF; x += delta {
				binary.ReadUvarint(r)
			}
			f.Close()
			os.Remove(file)
		case bigUvarintCase:
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				byteNum += binary.PutUvarint(buff, x)
				binary.Uvarint(buff)
			}
		case bigUvarintCaseRead:
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				byteNum += binary.PutUvarint(buff, x)
				reader := binary.BytesReader(buff)
				binary.ReadUvarint(&reader)
			}
		case bigUvarintCaseReadFile:
			file := "uvarintbig.hex"
			f, _ := os.Create(file)
			w := bufio.NewWriter(f)
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				n := binary.PutUvarint(buff, x)
				byteNum += n
				w.Write(buff[:n])
			}
			f.Sync()
			f.Close()
			f, _ = os.Open(file)
			r := bufio.NewReader(f)
			for x := uint64(deltaBig); x != 0; x += deltaBig {
				binary.ReadUvarint(r)
			}
			f.Close()
			os.Remove(file)
		}

	case BenchDecode:
		switch datax := data.(type) {
		case []byte:
			for i := 0; i < doCnt; i++ {
				byteNum = 0
				for byteNum < len(datax) {
					_, n := binary.Uvarint(datax[byteNum:])
					byteNum += n
				}
			}
		case littleUvarintCase:
		case bigUvarintCase:
		}

	}
	dur := Duration(time.Now().Sub(start))
	speed = Speed(float64(time.Duration(byteNum)*time.Second) / float64(dur*1024*1024))
	return dur, speed, Size(byteNum)
}

// DoBench runs a bench test case for binary
func DoBench(bench benchType, data interface{},
	doCnt int, enableSerializer bool) (t Duration, speed Speed, size Size) {
	start := time.Now()
	var err error
	var b []byte
	byteNum := 0
	switch bench {
	case BenchStdWrite:
		s := std.Size(data)
		if s < 0 {
			return -1, Speed(s), 0
		}
		for i := 0; i < doCnt; i++ {
			buffer.Reset()
			std.Write(buffer, std.LittleEndian, data)
		}
		byteNum = s * doCnt
	case BenchStdRead:
		s := std.Size(data)
		if s < 0 {
			return -1, Speed(s), 0
		}
		if err = std.Write(buffer, std.LittleEndian, data); err != nil {
			panic(err)
		}
		w := NewSame(data)
		b := buffer.Bytes()
		for i := 0; i < doCnt; i++ {
			r := binary.BytesReader(b)
			std.Read(&r, std.LittleEndian, w)
		}
		byteNum = s * (doCnt + 1)
	case BenchEncode:
		for i := 0; i < doCnt; i++ {
			b, err = binary.EncodeX(data, buff, enableSerializer)
			byteNum += len(b)
		}
		if err != nil {
			panic(err)
		}
	case BenchDecode:
		b, err := binary.EncodeX(data, buff, enableSerializer)
		w := NewSame(data)
		for i := 0; i < doCnt; i++ {
			err = binary.DecodeX(b, w, enableSerializer)
		}
		byteNum = len(b) * (doCnt + 1)
		if err != nil {
			panic(err)
		}

	case BenchGobEncode:
		buffer.Reset()
		coder := gob.NewEncoder(buffer)
		if err := coder.Encode(data); err != nil {
			panic(err)
		}
		byteNum = buffer.Len() * (doCnt + 1)
		for i := 0; i < doCnt; i++ {
			buffer.Reset()
			coder.Encode(data)
		}

	case BenchGobDecode:
		buffer.Reset()
		encoder := gob.NewEncoder(buffer)
		if err := encoder.Encode(data); err != nil {
			panic(err)
		}
		buf := buffer.Bytes()
		byteNum = buffer.Len() * (doCnt + 1)
		r := binary.BytesReader(buf)
		decoder := gob.NewDecoder(&r)
		w := NewSame(data)
		decoder.Decode(w)
		for i := 0; i < doCnt; i++ {
			r = binary.BytesReader(buf)
			decoder.Decode(w)
		}
	}

	dur := Duration(time.Now().Sub(start))
	//speed = Speed(float64(time.Duration(byteNum)*time.Second) / float64(dur*1024*1024))
	speed = Speed(float64(time.Duration(doCnt)*time.Second) / float64(dur))
	return dur, speed, Size(byteNum)
}

func NewSame(x interface{}) interface{} {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_new := reflect.New(t)
	//	switch t.Kind() {
	//	case reflect.Slice:
	//		_new.Set(reflect.MakeSlice(t, t.Len(), t.Len()))
	//	case reflect.Map:
	//		_new.Set(reflect.MakeMap(t))
	//	}
	return _new.Interface()
}

type (
	Bool       bool
	Int8       int8
	Int16      int16
	Int32      int32
	Int64      int64
	Uint8      uint8
	Uint16     uint16
	Uint32     uint32
	Uint64     uint64
	Int        int
	Uint       uint
	Float32    float32
	Float64    float64
	Complex64  complex64
	Complex128 complex128
	String     string
)

type BaseStruct struct {
	Bool       Bool
	Int8       Int8
	Int16      Int16
	Int32      Int32
	Int64      Int64
	Uint8      Uint8
	Uint16     Uint16
	Uint32     Uint32
	Uint64     Uint64
	Float32    Float32
	Float64    Float64
	Complex64  Complex64
	Complex128 Complex128
}

type FastValues struct {
	Int             int
	Uint            uint
	Bool            bool
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	String          string
	IntSlice        []int
	UintSlice       []uint
	BoolSlice       []bool
	Int8Slice       []int8
	Int16Slice      []int16
	Int32Slice      []int32
	Int64Slice      []int64
	Uint8Slice      []uint8
	Uint16Slice     []uint16
	Uint32Slice     []uint32
	Uint64Slice     []uint64
	Float32Slice    []float32
	Float64Slice    []float64
	Complex64Slice  []complex64
	Complex128Slice []complex128
	StringSlice     []string
}

type NormalValues struct {
	Int             Int
	Uint            Uint
	Bool            Bool
	Int8            Int8
	Int16           Int16
	Int32           Int32
	Int64           Int64
	Uint8           Uint8
	Uint16          Uint16
	Uint32          Uint32
	Uint64          Uint64
	Float32         Float32
	Float64         Float64
	Complex64       Complex64
	Complex128      Complex128
	String          String
	IntSlice        []Int
	UintSlice       []Uint
	BoolSlice       []Bool
	Int8Slice       []Int8
	Int16Slice      []Int16
	Int32Slice      []Int32
	Int64Slice      []Int64
	Uint8Slice      []Uint8
	Uint16Slice     []Uint16
	Uint32Slice     []Uint32
	Uint64Slice     []Uint64
	Float32Slice    []Float32
	Float64Slice    []Float64
	Complex64Slice  []Complex64
	Complex128Slice []Complex128
	StringSlice     []String
	MapUU           map[uint64]uint32
}

type LargeData struct {
	LgIntSlice        []int
	LgUintSlice       []uint
	LgBoolSlice       []bool
	LgInt8Slice       []int8
	LgInt16Slice      []int16
	LgInt32Slice      []int32
	LgInt64Slice      []int64
	LgUint8Slice      []uint8
	LgUint16Slice     []uint16
	LgUint32Slice     []uint32
	LgUint64Slice     []uint64
	LgFloat32Slice    []float32
	LgFloat64Slice    []float64
	LgComplex64Slice  []complex64
	LgComplex128Slice []complex128
	LgStringSlice     []string
	LgIntArray        [1000]int
	LgUintArray       [1000]uint
	LgBoolArray       [1000]bool
	LgInt8Array       [1000]int8
	LgInt16Array      [1000]int16
	LgInt32Array      [1000]int32
	LgInt64Array      [1000]int64
	LgUint8Array      [1000]uint8
	LgUint16Array     [1000]uint16
	LgUint32Array     [1000]uint32
	LgUint64Array     [1000]uint64
	LgFloat32Array    [1000]float32
	LgFloat64Array    [1000]float64
	LgComplex64Array  [1000]complex64
	LgComplex128Array [1000]complex128
	LgStringArray     [1000]string
	LgMapUU           map[uint64]uint32
}

type FullStruct struct {
	FastValues   FastValues
	NormalValues NormalValues
	LargeData    LargeData
	Special      Special
}

type Special struct {
	RegedStruct RegedStruct
	Serializer  Serializer
	BaseStruct  BaseStruct
}

type RegedStruct BaseStruct

type Serializer BaseStruct

func (s Serializer) Size() int {
	return binary.SizeX(BaseStruct(s), false)
}

func (s Serializer) Encode(buffer []byte) ([]byte, error) {
	return binary.Encode(BaseStruct(s), buffer)
}

func (s *Serializer) Decode(buffer []byte) error {
	return binary.Decode(buffer, (*BaseStruct)(s))
}

type BenchCase struct {
	ID     int
	Name   string
	Length int
	Data   interface{}
}

var cases []*BenchCase

func init() {
	binary.RegisterType((*RegedStruct)(nil))
	binary.RegisterType((*Serializer)(nil))

	rnd := random.NewRand(0)
	seed := uint64(1)
	rnd.ValueX(&full.FastValues, seed, 10, 0)
	rnd.ValueX(&full.NormalValues, seed, 10, 0)
	rnd.ValueX(&full.LargeData, seed, 1000, 1000)
	rnd.ValueX(&full.Special, seed, 10, 0)
	//fmt.Printf("%@#v\n", full)

	genCase(full.FastValues, "FastValues")
	genCase(full.NormalValues, "NormalValues")
	genCase(full.LargeData, "LargeData")
	genCase(full.Special, "Special")
	//fmt.Printf("%@#v\n", cases)
}

func genCase(data interface{}, name string) {
	v := reflect.ValueOf(data)
	t := v.Type()
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		finfo := t.Field(i)
		c := &BenchCase{}
		c.ID = len(cases) + 1
		c.Name = name + "." + finfo.Name
		f := v.Field(i)
		c.Data = f.Interface()
		switch k := f.Kind(); k {
		case reflect.Slice, reflect.Array, reflect.Map:
			c.Length = f.Len()
		case reflect.Struct:
			c.Length = f.NumField()
		}
		if c.Length == 0 {
			c.Length = 1
		}
		cases = append(cases, c)
	}
}
