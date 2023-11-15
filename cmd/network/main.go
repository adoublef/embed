package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adoublef/mvp/cmd/network/server"
	eg "github.com/adoublef/mvp/errgroup"
	"github.com/adoublef/mvp/nats"
	sql "github.com/adoublef/mvp/sqlite3"
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
	db, err := sql.Open(o.DSN)
	if err != nil {
		return err
	}
	defer db.Close()
	
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

	s, err := server.NewServer(o.Addr, nc, db)
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
	DSN string
}

// Parse parses flags definitions and runtime environment variables
func parse() (*Options, error) {
	fs, args := flag.NewFlagSet("", flag.ExitOnError), os.Args[1:]
	o := &Options{}

	fs.StringVar(&o.Addr, "addr", ":8080", "http listen address")
	fs.StringVar(&o.DSN, "dsn", "main.db", "database path")
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	return o, nil
}