package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/adoublef/embed/cmd/whiteboard/server"
	eg "github.com/adoublef/embed/errgroup"
	"github.com/adoublef/embed/nats"
	sql "github.com/adoublef/embed/sqlite3"
	"github.com/choria-io/fisk"
)

type serveCommand struct {
	addr string
	js string
	dsn  string
}

func configureServeCommand(app commandHost) {
	c := &serveCommand{}
	serve := app.Command("serve", "Run application server").Alias("s").Action(c.serve)
	serve.Flag("addr", "Listen address").StringVar(&c.addr)
	serve.Flag("js", "JetStream directory").StringVar(&c.js)
	serve.Flag("dsn", "Datasource name").StringVar(&c.dsn)
}

func init() {
	registerCommand("serve", 0, configureServeCommand)
}

func (c *serveCommand) serve(_ *fisk.ParseContext) error {
	if c.addr == "" {
		c.addr = ":8080"
	}
	if c.dsn == "" {
		c.dsn = ":memory:"
	}
	if c.js == "" {
		c.js = "."
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := sql.Open(c.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	ns, err := nats.NewServer(c.js)
	if err != nil {
		return err
	}
	ns.Wait()
	nc, err := nats.Connect(ns)
	if err != nil {
		return err
	}
	defer nc.Close()

	s, err := server.NewServer(c.addr, nc, db)
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
