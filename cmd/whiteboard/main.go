package main

import (
	"os"
	"sort"
	"sync"

	"github.com/adoublef/embed/log"
	"github.com/choria-io/fisk"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ncli := fisk.New("", "__EMPTY__")
	// use charm for a pretty print
	ncli.UsageWriter(os.Stdout)
	ncli.HelpFlag.Short('h')
	ncli.WithCheats().CheatCommand.Hidden()

	sort.Slice(commands, func(i int, j int) bool {
		return commands[i].Name < commands[j].Name
	})
	for _, c := range commands {
		c.Command(ncli)
	}
	
	_, err := ncli.Parse(os.Args[1:])
	return err
}

type command struct {
	Name    string
	Order   int
	Command func(app commandHost)
}

type commandHost interface {
	Command(name string, help string) *fisk.CmdClause
}

var (
	commands = []*command{}
	mu       sync.Mutex
)

func registerCommand(name string, order int, c func(app commandHost)) {
	mu.Lock()
	commands = append(commands, &command{name, order, c})
	mu.Unlock()
}
