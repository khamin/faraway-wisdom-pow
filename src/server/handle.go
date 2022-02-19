package server

import (
	"encoding/binary"
	"net"

	"github.com/sirupsen/logrus"
)

const N, K uint32 = 102, 5

var BytesOrder = binary.LittleEndian

var read = binary.Read
var write = binary.Write

// Handle client request.
func (srv *Server) handler(conn net.Conn) {
	logFields := logrus.WithField("remoteAddr", conn.RemoteAddr())
	logFields.Info("connected")

	defer conn.Close()
	defer logFields.Info("disconnected")

	ok, err := srv.challenge(conn)

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
