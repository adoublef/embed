package nats

import (
	"fmt"

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
		return nil, fmt.Errorf("connect: %w", err)
	}
	
	return &Conn{nc: nc}, nil
}

// JetStream returns a JetStreamContext for messaging and stream management. 
func JetStream(nc *Conn) (JetStreamContext, error) {
	jsc, err := nc.nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("jet stream context: %w", err)
	}
	return jsc, nil
}