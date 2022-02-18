package equihash

import (
	"testing"
)

func TestEquihash(t *testing.T) {
	var n, k uint32 = 102, 5

	seed := NewSeed()

	config := Config{
		K:    k,
		N:    n,
		Seed: seed,
	}

	hash := Equihash{
		Config: config,
	}

	proof := hash.FindProof()

	if len(proof.Inputs) == 0 {
		t.Error("solutions not found")
	}

	if !proof.Test() {
		t.Error("proof test failed")
	}
}
