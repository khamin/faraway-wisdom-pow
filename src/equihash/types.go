package equihash

const SeedLen = 4

type Config struct {
	K     uint32
	N     uint32
	Nonce uint32
	Seed  Seed
}

type Tuple struct {
	blocks []uint32
	ref    uint32
}

type Tuples []Tuple

type Fork struct {
	Ref1 uint32
	Ref2 uint32
}

type Forks []Fork

type Seed [SeedLen]uint32
