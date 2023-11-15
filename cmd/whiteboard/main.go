package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/adoublef/mvp/cmd/whiteboard/server"
	eg "github.com/adoublef/mvp/errgroup"
	"github.com/adoublef/mvp/nats"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// TODO parse flags

	if err := run(ctx); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context) error {
	ns, err := nats.NewServer()
	if err != nil {
		return err
	}
	ns.Wait()
	// add a nats connection
	nc, err := nats.Connect(ns)
	if err != nil {
		return err
	}
	defer nc.Close()

	s, err := server.NewServer(nc)
	if err != nil {
		return err
	}
	
	g := eg.New(ctx)
	g.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return s.Shutdown()
	})
	return g.Wait()
}
