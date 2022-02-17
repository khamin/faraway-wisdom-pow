package server

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	config   *Config
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
}

// Start server.
func (srv *Server) start() {
	defer srv.wg.Done()

	for {
		conn, err := srv.listener.Accept()

		if err != nil {
			select {
			case <-srv.quit:
				return
			default:
				logrus.WithField("err", err).Error("accept error")
				continue
			}
		}

		srv.wg.Add(1)

		go func() {
			defer srv.wg.Done()
			srv.handler(conn)
		}()
	}
}

// Run server shutdown routine.
// Return error if any.
func (srv *Server) Stop() error {
	close(srv.quit)
	srv.listener.Close()

	c := make(chan struct{})

	go func() {
		defer close(c)
		srv.wg.Wait()
	}()

	select {
	case <-c:
		return nil
	case <-time.After(srv.config.StopTimeout):
		return fmt.Errorf("timeout")
	}
}
