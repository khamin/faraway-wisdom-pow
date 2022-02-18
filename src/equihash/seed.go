package equihash

import "math/rand"

type Seed [SeedLen]uint32

func NewSeed() (seed Seed) {
	for i := 0; i < len(seed); i++ {
		seed[i] = rand.Uint32()
	}

	return
}
