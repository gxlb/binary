package bench

import (
	"testing"
)

func TestRand(t *testing.T) {
	rand.Srand(RandSeed32(0))
	println(rand.String(12))
	println(rand.Float64())
	println(rand.Complex128())
	println(rand.Uint32())
	println(rand.Uint64())
	println(rand.Float32())
}
