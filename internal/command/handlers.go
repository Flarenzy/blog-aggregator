package command

import (
	"context"
	"errors"
	"fmt"
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

func handlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func handlerUsers(s *State, _ Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %v (current)\n", user.Name)
		} else {
			fmt.Printf("* %v\n", user.Name)
		}
	}
	return nil
}

func handlerAgg(s *State, cmd Command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	xml, err := unescapeXML(rss)
	if err != nil {
		return err
	}
	fmt.Println(xml)
	return nil
}
