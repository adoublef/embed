package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/adoublef/mvp/cmd/whiteboard/server"
	"github.com/adoublef/mvp/errgroup"
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
	s, err := server.NewServer()
	if err != nil {
		return err
	}
	g := errgroup.New(ctx)
	g.Go(func(ctx context.Context) error {
		return s.ListenAndServe()
	})
	g.Go(func(ctx context.Context) error {
		<-ctx.Done()
		return s.Shutdown()
	})
	return g.Wait()
}
