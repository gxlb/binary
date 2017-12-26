//bench test case for binary and std

package bench

import (
	"bytes"
	std "encoding/binary"
	"reflect"
	"time"

	"github.com/vipally/binary"
)

type BenchCase struct {
	Id               int
	Name             string
	DoCnt            int
	EnableSerializer bool
	Data             interface{}
}

var cases []*BenchCase

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

// DoBench runs a bench test case for binary
func DoBench(bench benchType, data interface{}, doCnt int, enableSerializer bool, name string) time.Duration {
	start := time.Now()
	var err error
	switch bench {
	case BenchStdWrite:
		s := std.Size(data)
		if s <= 0 {
			println(name, "unsupported ")
			return 0
		}
		for i := 0; i < doCnt; i++ {
			buffer.Reset()
			std.Write(buffer, std.LittleEndian, data)
		}
	case BenchStdRead:
		s := std.Size(data)
		if s <= 0 {
			println(name, "unsupported ")
			return 0
		}
		if err = std.Write(buffer, std.LittleEndian, data); err != nil {
			panic(err)
		}
		w := newSame(data)
		b := buffer.Bytes()
		for i := 0; i < doCnt; i++ {
			r := binary.BytesReader(b)
			std.Read(&r, std.LittleEndian, w)
		}
	case BenchEncode:
		for i := 0; i < doCnt; i++ {
			_, err = binary.EncodeX(data, buff, enableSerializer)
		}
		if err != nil {
			panic(err)
		}
	case BenchDecode:
		std.Write(buffer, std.LittleEndian, data)
		w := newSame(data)
		b := buffer.Bytes()
		for i := 0; i < doCnt; i++ {
			err = binary.DecodeX(b, w, enableSerializer)
		}
		if err != nil {
			panic(err)
		}
	}

	dur := time.Now().Sub(start)
	return dur
}

func newSame(x interface{}) interface{} {
	t := reflect.TypeOf(x)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	_new := reflect.New(t)
	switch t.Kind() {
	case reflect.Slice:
		_new.Set(reflect.MakeSlice(t, t.Len(), t.Len()))
	case reflect.Map:
		_new.Set(reflect.MakeMap(t))
	}
	return _new.Interface()
}

func init() {
	cases = []*BenchCase{
		&BenchCase{0, "int", benchDoCnt, false, int(0)},
	}
}
