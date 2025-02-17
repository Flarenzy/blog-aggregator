package command

import (
	"context"
	"github.com/Flarenzy/blog-aggregator/internal/database"
)

func middlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(s *State, cmd Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}
}
