package main

import "github.com/choria-io/fisk"

type migrateCommand struct {
	dsn string
}

func configureMigrateCommand(app commandHost) {
	c := &migrateCommand{}
	migrate := app.Command("migrate", "Run database migrations").Alias("m").Action(c.migrate)
	migrate.Flag("dsn", "Datasource name").PlaceHolder("DSN").StringVar(&c.dsn)
}

func init() {
	registerCommand("migrate", 1, configureMigrateCommand)
}

func (c *migrateCommand) migrate(_ *fisk.ParseContext) error {
	if c.dsn == "" {
		c.dsn = ":memory:"
	}
	// run migrations
	return nil
}