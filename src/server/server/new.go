package server

import (
	"net"
	"time"
	"wisdom/server/quotes"
)

// Represents server config.
type Config struct {
	Addr        string
	Quotes      *quotes.Quotes
	StopTimeout time.Duration
}

// Start and return new server instance.
func New(config *Config) (*Server, error) {
	listener, err := net.Listen("tcp", config.Addr)

	if err != nil {
		return nil, err
	}

	quit := make(chan interface{})

	srv := &Server{
		config:   config,
		listener: listener,
		quit:     quit,
	}

	srv.wg.Add(1)
	go srv.start()

	return srv, nil
}
