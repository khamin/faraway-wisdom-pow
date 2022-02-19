package equihash

import "math/rand"

// Return new instance of equihash.
// Nil seed will be generated automatically.
func New(n, k uint32, seed *Seed) *Equihash {
	if seed == nil {
		seed = &Seed{}

		for i := 0; i < len(seed); i++ {
			seed[i] = rand.Uint32()
		}
	}

	config := Config{
		N:    n,
		K:    k,
		Seed: *seed,
	}

	return &Equihash{
		Config: config,
	}
}
