package server

import (
	"encoding/binary"
	"net"
	"time"

	"wisdom/equihash"

	"github.com/sirupsen/logrus"
)

const N, K uint32 = 102, 5

var BytesOrder = binary.LittleEndian

var read = binary.Read
var write = binary.Write

// Perform equihash proof-of-work challenge.
func challenge(conn net.Conn) (bool, error) {
	logFields := logrus.WithField("remoteAddr", conn.RemoteAddr())

	deadline := time.Now().Add(10 * time.Second)

	if err := conn.SetDeadline(deadline); err != nil {
		return false, err
	}

	defer conn.SetDeadline(time.Time{})

	if err := write(conn, BytesOrder, N); err != nil {
		return false, err
	}

	if err := write(conn, BytesOrder, K); err != nil {
		return false, err
	}

	seed := equihash.NewSeed()

	if err := write(conn, BytesOrder, seed); err != nil {
		return false, err
	}

	logFields.WithFields(logrus.Fields{
		"k":    K,
		"n":    N,
		"seed": seed,
	}).Info("challenge sent")

	var nonce uint32

	if err := read(conn, BytesOrder, &nonce); err != nil {
		return false, err
	}

	var inputs [32]uint32

	if err := read(conn, BytesOrder, &inputs); err != nil {
		return false, err
	}

	config := equihash.Config{
		K:     K,
		N:     N,
		Nonce: nonce,
		Seed:  seed,
	}

	proof := equihash.Proof{
		Config: config,
		Inputs: inputs[:],
	}

	start := time.Now()
	ok := proof.Test()

	logFields.WithFields(logrus.Fields{
		"inputs": proof.Inputs,
		"nonce":  proof.Config.Nonce,
		"seed":   proof.Config.Seed,
		"time":   time.Since(start),
	}).Debug("proof test")

	return ok, nil
}

// Handle client request.
func (srv *Server) handler(conn net.Conn) {
	logFields := logrus.WithField("remoteAddr", conn.RemoteAddr())
	logFields.Info("connected")

	defer conn.Close()
	defer logFields.Info("disconnected")

	ok, err := challenge(conn)

	if !ok {
		logFields.Error("auth failed")
		return
	}

	if err != nil {
		logFields.WithField("err", err).Error("auth error")
		return
	}

	logFields.Debug("challenge passed")
	quote := srv.config.Quotes.Pick()

	if err := write(conn, BytesOrder, quote.Format()); err != nil {
		logFields.WithField("err", err).Error("write failed")
		return
	}
}
