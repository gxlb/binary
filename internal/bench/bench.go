//bench test case for binary and std

package bench

import (
	"bytes"
	std "encoding/binary"
	"fmt"
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
	return fmt.Sprintf("%.02fMB/s", s)
}

// DoBench runs a bench test case for binary
func DoBench(bench benchType, data interface{},
	doCnt int, enableSerializer bool, name string) (t time.Duration, speed Speed) {
	start := time.Now()
	var err error
	var b []byte
	byteNum := 0
	switch bench {
	case BenchStdWrite:
		s := std.Size(data)
		if s < 0 {
			return 0, Speed(s)
		}
		for i := 0; i < doCnt; i++ {
			buffer.Reset()
			std.Write(buffer, std.LittleEndian, data)
		}
		byteNum = s * doCnt
	case BenchStdRead:
		s := std.Size(data)
		if s < 0 {
			return 0, Speed(s)
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
	}

	dur := time.Now().Sub(start)
	speed = Speed(float64(time.Duration(byteNum)*time.Second) / float64(dur*1024*1024))
	return dur, speed
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
	Int        Int
	Uint       Uint
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

type NoneFastValues struct {
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
	FastValues     FastValues
	NoneFastValues NoneFastValues
	LargeData      LargeData
	Special        Special
}

type Special struct {
	RegedStruct RegedStruct
	Serializer  Serializer
	BaseStruct  BaseStruct
}

type RegedStruct BaseStruct

type Serializer BaseStruct

func (s Serializer) Size() int {
	return 0
}

func (s Serializer) Encode(buffer []byte) ([]byte, error) {
	return nil, nil
}

func (s *Serializer) Decode(buffer []byte) error {
	return nil
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
	rnd.ValueX(&full.NoneFastValues, seed, 10, 0)
	rnd.ValueX(&full.LargeData, seed, 1000, 1000)
	rnd.ValueX(&full.Special, seed, 10, 0)
	//fmt.Printf("%@#v\n", full)

	genCase(full.FastValues, "FastValues")
	genCase(full.NoneFastValues, "NoneFastValues")
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
		if k := f.Kind(); k == reflect.Slice || k == reflect.Array || k == reflect.Map {
			c.Length = f.Len()
		}
		if c.Length == 0 {
			c.Length = 1
		}
		cases = append(cases, c)
	}
}
