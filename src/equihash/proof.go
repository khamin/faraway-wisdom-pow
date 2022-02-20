package equihash

import "golang.org/x/crypto/blake2b"

type Proof struct {
	Config Config
	Nonce  uint32
	Inputs []uint32
}

func (p *Proof) Test() bool {
	input := make([]uint32, SeedLen+2)

	for i := uint32(0); i < SeedLen; i++ {
		input[i] = p.Config.Seed[i]
	}

	input[SeedLen] = p.Nonce
	input[SeedLen+1] = 0

	blocks := make([]uint32, p.Config.K+1)

	for i := uint32(0); i < uint32(len(p.Inputs)); i++ {
		input[SeedLen+1] = p.Inputs[i]

		hash := blake2b.Sum256(toBytes(input[:]))
		buf := fromBytes(hash[:])

		for j := uint32(0); j < (p.Config.K + 1); j++ {
			// Select j-th block of n/(k+1) bits.
			blocks[j] ^= buf[j] >> (32 - p.Config.N/(p.Config.K+1))
		}
	}

	b := true

	for j := uint32(0); j < (p.Config.K + 1); j++ {
		b = b && (blocks[j] == 0)
	}

	return b
}
