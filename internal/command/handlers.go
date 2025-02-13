package command

import (
	"context"
	"errors"
	"github.com/Flarenzy/blog-aggregator/internal/database"
	"github.com/google/uuid"
	"time"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the login handler expects a single argument, the username")
	}
	user, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	err = s.Config.SetUser(user.Name)
	if err != nil {
		return err
	}
	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("the register handler expects a single argument, the username")
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	_, err := s.Db.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	s.Logger.Info("user has been created with params",
		"Name", cmd.Args[0],
		"ID", userParams.ID,
		"CreateAt", userParams.CreatedAt,
		"UpdatedAt", userParams.UpdatedAt)
	return nil
}
