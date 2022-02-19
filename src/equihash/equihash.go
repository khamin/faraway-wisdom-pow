package equihash

import (
	"sort"

	"golang.org/x/crypto/blake2b"
)

const maxNonce = 0xFFFFF
const listLen = 5

// Maximum collision factor.
const forkMultiplier = 3

type Equihash struct {
	Config Config
	filled []uint32
	forks  []Forks
	sols   []Proof
	tuples []Tuples
}

func (e *Equihash) FindProof() Proof {
	e.Config.Nonce = 1

	for {
		if e.Config.Nonce >= maxNonce {
			break
		}

		e.Config.Nonce++

		e.init()
		e.fill()

		for i := uint32(1); i <= e.Config.K; i++ {
			// XOR collisions, concatenate indices and shift.
			e.resolveCollisions(i == e.Config.K)
		}

		// Duplicate check.
		for i := uint32(0); i < uint32(len(e.sols)); i++ {
			inputs := e.sols[i].Inputs

			sort.Slice(inputs, func(i, j int) bool {
				return inputs[i] < inputs[j]
			})

			dup := false

			for k := uint32(0); k < uint32(len(inputs)-1); k++ {
				if inputs[k] == inputs[k+1] {
					dup = true
				}
			}

			if !dup {
				return e.sols[i]
			}
		}
	}

	return Proof{
		Config: e.Config,
		Inputs: make([]uint32, 0),
	}
}

func (e *Equihash) fill() {
	length := uint32(4) << (e.Config.N/(e.Config.K+1) - 1)

	var input [SeedLen + 2]uint32

	for i := uint32(0); i < SeedLen; i++ {
		input[i] = e.Config.Seed[i]
	}

	input[SeedLen] = e.Config.Nonce
	input[SeedLen+1] = 0

	for i := uint32(0); i < length; i++ {
		hash := blake2b.Sum256(toBytes(input[:]))
		buf := fromBytes(hash[:])

		var index uint32 = buf[0] >> (32 - e.Config.N/(e.Config.K+1))
		var count uint32 = e.filled[index]

		if count < listLen {
			for j := uint32(1); j < e.Config.K+1; j++ {
				// Select j-th block of n/(k+1) bits.
				v := buf[j] >> (32 - e.Config.N/(e.Config.K+1))
				e.tuples[index][count].blocks[j-1] = v
			}

			e.tuples[index][count].ref = i
			e.filled[index] += 1
		}

		input[SeedLen+1]++
	}
}

func (e *Equihash) init() {
	length := uint32(1) << (e.Config.N / (e.Config.K + 1))
	e.tuples = make([]Tuples, length)

	for i := uint32(0); i < length; i++ {
		e.tuples[i] = make(Tuples, listLen)

		for j := 0; j < listLen; j++ {
			// k blocks to store, one left for index.
			e.tuples[i][j] = Tuple{
				blocks: make([]uint32, e.Config.K),
			}
		}
	}

	e.filled = make([]uint32, length)

	e.sols = nil
	e.forks = nil
}

func (e *Equihash) resolveCollisions(store bool) {
	// Number of rows in the hashtable
	tableLen := uint32(len(e.tuples))

	// Max number of collisions to be found.
	maxCollisions := uint32(len(e.tuples) * forkMultiplier)

	// Number of blocks in the future collisions.
	blocks := uint32(len(e.tuples[0][0].blocks) - 1)

	// List of forks created at this step.
	forks := make(Forks, maxCollisions)

	collisions := make([]Tuples, tableLen)

	for i := uint32(0); i < tableLen; i++ {
		collisions[i] = make(Tuples, listLen)

		for j := 0; j < listLen; j++ {
			collisions[i][j] = Tuple{
				blocks: make([]uint32, blocks),
			}
		}
	}

	// Number of entries in rows.
	filled := make([]uint32, tableLen)

	var collisionsCount uint32

	for i := uint32(0); i < tableLen; i++ {
		for j := uint32(0); j < e.filled[i]; j++ {
			for m := j + 1; m < e.filled[i]; m++ { // Collision.
				idx := e.tuples[i][j].blocks[0] ^ e.tuples[i][m].blocks[0]

				fork := Fork{
					Ref1: e.tuples[i][j].ref,
					Ref2: e.tuples[i][m].ref,
				}

				// Check if we get a solution.
				if store { // Last step.
					if idx == 0 { // Solution.
						inputs := e.resolveTree(fork)

						e.sols = append(e.sols, Proof{
							Config: e.Config,
							Inputs: inputs,
						})
					}
				} else { // Resolve.
					if filled[idx] < listLen && collisionsCount < maxCollisions {
						for l := uint32(0); l < blocks; l++ {
							collisions[idx][filled[idx]].blocks[l] = e.tuples[i][j].blocks[l+1] ^ e.tuples[i][m].blocks[l+1]
						}

						forks[collisionsCount] = fork
						collisions[idx][filled[idx]].ref = collisionsCount
						filled[idx]++
						collisionsCount++
					}
				}
			}
		}
	}

	e.forks = append(e.forks, forks)

	e.tuples = collisions
	e.filled = filled
}

func (e *Equihash) resolveTree(fork Fork) []uint32 {
	return e.resolveTreeByLevel(fork, uint32(len(e.forks)))
}

func (e *Equihash) resolveTreeByLevel(fork Fork, level uint32) []uint32 {
	if level == 0 {
		return []uint32{
			fork.Ref1,
			fork.Ref2,
		}
	}

	v1 := e.resolveTreeByLevel(e.forks[level-1][fork.Ref1], level-1)
	v2 := e.resolveTreeByLevel(e.forks[level-1][fork.Ref2], level-1)

	return append(v1, v2...)
}
