package command

import (
	"errors"
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("user has been set: %v\n", cmd.Args[0])
	return nil
}
