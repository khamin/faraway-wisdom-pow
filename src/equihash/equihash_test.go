package equihash

import (
	"testing"
)

const n, k = 102, 5
const nonce = 6

var seed = Seed{
	2596996162, 4039455774,
	2854263694, 1879968118,
}

var config = Config{
	K:    k,
	N:    n,
	Seed: seed,
}

var inputs = []uint32{
	9892, 14056, 20114, 29144,
	33169, 55503, 59828, 60514,
	69209, 70006, 75816, 97554,
	98161, 107885, 111381, 125631,
	127879, 135564, 139747, 148672,
	156601, 174964, 198971, 208584,
	212224, 214961, 224039, 226498,
	229762, 231585, 232063, 232488,
}

func TestEquihash(t *testing.T) {
	e := New(n, k, nil)
	proof := e.FindProof()

	if len(proof.Inputs) == 0 {
		t.Error("solutions not found")
	}

	if !proof.Test() {
		t.Error("proof test failed")
	}
}

func BenchmarkFindProof(b *testing.B) {
	e := Equihash{
		Config: config,
		Nonce:  nonce,
	}

	for i := 0; i < b.N; i++ {
		e.FindProof()
	}
}

func BenchmarkTest(b *testing.B) {
	proof := Proof{
		Config: config,
		Inputs: inputs,
		Nonce:  nonce,
	}

	for i := 0; i < b.N; i++ {
		if !proof.Test() {
			b.Fatal("test failed")
		}
	}
}
