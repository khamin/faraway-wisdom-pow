package equihash

import (
	"testing"
)

func TestEquihash(t *testing.T) {
	e := New(102, 5, nil)
	proof := e.FindProof()

	if len(proof.Inputs) == 0 {
		t.Error("solutions not found")
	}

	if !proof.Test() {
		t.Error("proof test failed")
	}
}
