package command

import (
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"log/slog"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	Config *config.Config
}

func NewState(c *config.Config) *State {
	return &State{c}
}

type Commands struct {
	registered map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	cmds := &Commands{make(map[string]func(*State, Command) error)}
	cmds.register("login", handlerLogin)
	return cmds
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	_, ok := c.registered[name]
	if !ok {
		c.registered[name] = f
		slog.Info("registering command", "name", name)
	}
}

func (c *Commands) Run(s *State, cmd Command) error {
	f := c.registered[cmd.Name]
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
