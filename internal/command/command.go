package command

import (
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"github.com/Flarenzy/blog-aggregator/internal/database"
	"log/slog"
	"strings"
)

type Command struct {
	Name string
	Args []string
}

type State struct {
	Config *config.Config
	Db     *database.Queries
	Logger *slog.Logger
}

func NewState(c *config.Config, db *database.Queries, logger *slog.Logger) *State {
	return &State{c, db, logger}
}

type Commands struct {
	registered map[string]func(*State, Command) error
}

func NewCommands() *Commands {
	cmds := &Commands{make(map[string]func(*State, Command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
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
	f, ok := c.registered[cmd.Name]
	if !ok {
		return fmt.Errorf("command %s not found", cmd.Name)
	}
	s.Logger.Debug("Running command", "name", cmd.Name, "args", strings.Join(cmd.Args, " "))
	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}
