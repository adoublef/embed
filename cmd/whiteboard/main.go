package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adoublef/mvp/cmd/whiteboard/server"
	eg "github.com/adoublef/mvp/errgroup"
	"github.com/adoublef/mvp/nats"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	o, err := parse()
	if err != nil {
		log.Fatalln(err)
	}

	if err := run(ctx, o); err != nil {
		log.Fatalln(err)
	}
}

func run(ctx context.Context, o *Options) error {
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

	s, err := server.NewServer(o.Addr, nc)
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

type Options struct {
	Addr string
}

// Parse parses flags definitions and runtime environment variables
func parse() (*Options, error) {
	fs, args := flag.NewFlagSet("", flag.ExitOnError), os.Args[1:]
	o := &Options{}

	fs.StringVar(&o.Addr, "addr", ":8080", "http listen address")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return o, nil
}