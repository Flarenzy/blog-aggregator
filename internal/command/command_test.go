package command

import (
	"errors"
	"fmt"
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
		//if &f != &g {
		//	t.Errorf("command %s has registered handler, expected func adr %v, got %v", k, &g, &f)
		//	t.Fatal()
		//}
	}
}

func TestCommandsRun(t *testing.T) {
	cmds := NewCommands()
	c := tempConfig(t)
	s := NewState(c)
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
			t.Fatal()
		}
	}
}

func TestCommandsRunUnknownCommand(t *testing.T) {
	cmds := NewCommands()
	c := tempConfig(t)
	s := NewState(c)
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
			t.Fatal()
		}
	}
}
