package nats

import (
	"log"

	"github.com/cenkalti/backoff"
	"github.com/nats-io/nats-server/v2/server"
)

type Server struct {
	ns *server.Server
}

// Wait will block until the server is ready for connections.
func (s *Server) Wait()  {
	b := backoff.NewExponentialBackOff()

	for {
		d := b.NextBackOff()
		ready := s.ns.ReadyForConnections(d)
		if ready {
			break
		}

		log.Printf("NATS server not ready, waited %s, retrying...", d)
	}
}

// NewServer will setup a new embedded nats server.
func NewServer() (*Server, error) {
	ns, err := server.NewServer(&server.Options{})
	if err != nil {
		return nil, err
	}
	ns.Start()
	return &Server{ns}, nil
}