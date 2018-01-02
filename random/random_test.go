package random

import (
	//"fmt"
	"testing"
)

type all struct {
	String      string
	PString     *string
	Uint        uint
	StringSlice []string
	StringMap   map[uint16]string
	Uint8Slice  []uint8
	Uint32Slice []uint32
}

func TestRand(t *testing.T) {
	Default.Srand(RandSeed(0))
	//	var cnt [256]int
	//	total := 6553600
	//	for i := 0; i < total; i++ {
	//		r := Default.Uint16()
	//		if r < 256 {
	//			cnt[r]++
	//		}

	//	}
	//	for i, v := range cnt {
	//		fmt.Printf("%3d %4d/%4d\n", i, v, total/65536)
	//	}
	println(Default.String(12))
	println(Default.Float64())
	println(Default.Complex128())
	println(Default.Uint32())
	println(Default.Uint64())
	println(Default.Float32())
	for i := 0; i < 50; i++ {
		for j := 0; j < 8; j++ {
			print(Default.Uint8(), "\t")
		}
		println("")
	}
	//	var v all
	//	Default.Value(&v)
	//	fmt.Printf("%@v\n", v)

}
