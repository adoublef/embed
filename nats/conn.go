package nats

import (
	"github.com/nats-io/nats.go"
)

type Conn struct {
	nc *nats.Conn
}

func (c *Conn) Close() error {
	if c.nc != nil {
		c.nc.Close()
	}
	return nil
}

// Connect will attempt to connect to an embedded nats server.
func Connect(ns *Server) (*Conn, error) {
	nc, err := nats.Options{
		InProcessServer: ns.ns,
	}.Connect()
	if err != nil {
		return nil, err
	}
	
	return &Conn{nc: nc}, nil
}