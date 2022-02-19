package server

import (
	"net"
	"time"
	"wisdom/equihash"

	"github.com/sirupsen/logrus"
)

// Perform equihash proof-of-work challenge.
func (srv *Server) challenge(conn net.Conn) (bool, error) {
	logFields := logrus.WithField("remoteAddr", conn.RemoteAddr())

	deadline := time.Now().Add(srv.config.StopTimeout)

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

	e := equihash.New(N, K, nil)

	if err := write(conn, BytesOrder, e.Config.Seed); err != nil {
		return false, err
	}

	logFields.WithFields(logrus.Fields{
		"k":    e.Config.K,
		"n":    e.Config.N,
		"seed": e.Config.Seed,
	}).Info("challenge sent")

	if err := read(conn, BytesOrder, &e.Config.Nonce); err != nil {
		return false, err
	}

	var inputs [32]uint32

	if err := read(conn, BytesOrder, &inputs); err != nil {
		return false, err
	}

	proof := equihash.Proof{
		Config: e.Config,
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
