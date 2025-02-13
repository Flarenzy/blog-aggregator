package command

import (
	"errors"
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/logging"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestNewCommands(t *testing.T) {
	cmds := NewCommands()
	knowCommands := []string{"login"}
	knowHandlers := make(map[string]func(s *State, c Command) error)
	knowHandlers["login"] = handlerLogin
	for _, k := range knowCommands {
		_, ok := cmds.registered[k]
		if !ok {
			t.Errorf("command %s not found", k)
			t.Fatal()
		}
		_, ok = knowHandlers[k]
		if !ok {
			t.Errorf("command %s doesn't have know handler", k)
			t.Fatal()
		}
	}
}

func TestCommandsRun(t *testing.T) {
	cmds := NewCommands()
	c := tempConfig(t)
	randNum := strconv.Itoa(rand.Int())
	logger, f, err := logging.NewLogger("gator-temp"+randNum+".log", slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}
	s := NewState(c, logger)
	cases := []struct {
		cmd     string
		handler func(s *State, c Command) error
		name    string
	}{
		{
			cmd:     "login",
			handler: handlerLogin,
			name:    "Flarenzy",
		},
	}
	for _, ca := range cases {
		err := cmds.Run(s,
			Command{
				Name: ca.cmd,
				Args: []string{ca.name},
			})
		if err != nil {
			t.Fatal(err)
		}
		if s.Config.CurrentUserName != ca.name {
			t.Errorf("command %s expected %s got %s", ca.cmd, ca.name, s.Config.CurrentUserName)
			f.Close()
			os.Remove("gator-temp" + randNum + ".log")
			t.Fatal()
		}
	}
	f.Close()
	os.Remove(f.Name())
}

func TestCommandsRunUnknownCommand(t *testing.T) {
	cmds := NewCommands()
	c := tempConfig(t)
	randNum := strconv.Itoa(rand.Int())
	logger, f, err := logging.NewLogger("gator-temp"+randNum+".log", slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}
	s := NewState(c, logger)
	cases := []struct {
		cmd     string
		handler func(s *State, c Command) error
		name    string
	}{
		{
			cmd:     "randomCommand",
			handler: func(s *State, c Command) error { return errors.New("command not found") },
			name:    "RandomCommand",
		},
	}
	for _, ca := range cases {
		err := cmds.Run(s,
			Command{
				Name: ca.cmd,
				Args: []string{ca.name},
			})
		expectedError := fmt.Sprintf("command %s not found", ca.cmd)
		if err != nil && err.Error() != expectedError {
			t.Errorf("command %s expected %s got %s", ca.cmd, expectedError, err.Error())
			f.Close()
			os.Remove("./gator-temp" + randNum + ".log")
			t.Fatal()
		}
	}
	f.Close()
	os.Remove(f.Name())
}
