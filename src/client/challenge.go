package client

import (
	"net"
	"time"
	"wisdom/equihash"

	"github.com/sirupsen/logrus"
)

func challenge(conn net.Conn) (bool, error) {
	var n, k uint32

	if err := read(conn, BytesOrder, &n); err != nil {
		return false, err
	}

	if err := read(conn, BytesOrder, &k); err != nil {
		return false, err
	}

	var seed [equihash.SeedLen]uint32

	if err := read(conn, BytesOrder, &seed); err != nil {
		return false, err
	}

	logrus.WithFields(logrus.Fields{
		"k":    k,
		"n":    n,
		"seed": seed,
	}).Debug("challenge received")

	config := equihash.Config{
		K:    k,
		N:    n,
		Seed: seed,
	}

	hash := equihash.Equihash{
		Config: config,
	}

	start := time.Now()
	proof := hash.FindProof()

	logrus.WithFields(logrus.Fields{
		"inputs": proof.Inputs,
		"nonce":  proof.Config.Nonce,
		"time":   time.Since(start),
	}).Debug("proof found")

	if err := write(conn, BytesOrder, proof.Config.Nonce); err != nil {
		return false, err
	}

	if err := write(conn, BytesOrder, proof.Inputs); err != nil {
		return false, err
	}

	return true, nil
}
