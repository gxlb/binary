package random

import (
	//"fmt"
	"testing"
)

func TestRand(t *testing.T) {
	rand.Srand(RandSeed(0))
	//	var cnt [256]int
	//	total := 6553600
	//	for i := 0; i < total; i++ {
	//		r := rand.Uint16()
	//		if r < 256 {
	//			cnt[r]++
	//		}

	//	}
	//	for i, v := range cnt {
	//		fmt.Printf("%3d %4d/%4d\n", i, v, total/65536)
	//	}
	println(rand.String(12))
	println(rand.Float64())
	println(rand.Complex128())
	println(rand.Uint32())
	println(rand.Uint64())
	println(rand.Float32())
	for i := 0; i < 50; i++ {
		for j := 0; j < 8; j++ {
			print(rand.Uint8(), "\t")
		}
		println("")
	}
}
