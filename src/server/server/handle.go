package server

import (
	"log"
	"net"

	"github.com/sirupsen/logrus"
)

// Handle client request.
func (srv *Server) handler(conn net.Conn) {
	remoteAddr := conn.RemoteAddr()

	logFields := logrus.WithField("remoteAddr", remoteAddr)
	logFields.Info("connected")

	defer conn.Close()

	quote := srv.config.Quotes.Pick()
	n, err := conn.Write(quote.Format())

	if err != nil {
		log.Print(err)
		return
	}

	logFields = logFields.WithField("bytesWritten", n)
	logFields.Info("disconnected")
}
