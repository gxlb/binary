package bench

import (
	"testing"
)

func TestRand(t *testing.T) {
	rand.Srand(RandSeed(0))
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
