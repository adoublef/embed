package nats

import (
	"fmt"

	"github.com/adoublef/embed/log"
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
func NewServer(js string) (*Server, error) {
	ns, err := server.NewServer(&server.Options{
		// ServerName: "",
		Port: 4222,
		// HTTPPort: 8222,
		JetStream: true,
		StoreDir: js,
	})
	if err != nil {
		return nil, fmt.Errorf("new nats server: %w", err)
	}
	ns.Start()
	return &Server{ns}, nil
}