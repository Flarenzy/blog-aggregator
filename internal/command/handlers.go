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
	if len(cmd.Args) <= 0 {
		s.Logger.Error("command expects the time between requests", "cmd", cmd.Name, "args", cmd.Args)
		return errors.New("the aggregate handler expects a single argument, the time between requests")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		s.Logger.Error("command expects the time between requests", "cmd", cmd.Name, "args", cmd.Args)
		return err
	}
	ticker := time.NewTicker(timeBetweenRequests)
	fmt.Printf("Collecting feeds every %v \n", cmd.Args[0])
	for ; ; <-ticker.C {
		fmt.Printf("Fetching all feeds for users %v\n", s.Config.CurrentUserName)
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("the add feed handler expects at least two arguments: name and url")
	}
	name := cmd.Args[0]
	url := cmd.Args[1]
	var feedParams database.CreateFeedParams
	curTime := time.Now()
	feedParams.ID = uuid.New()
	feedParams.Name = name
	feedParams.Url = url
	feedParams.UpdatedAt = curTime
	feedParams.CreatedAt = curTime
	feedParams.UserID = user.ID
	_, err := s.Db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		s.Logger.Error("failed to create feed", "url", url, "err", err)
		return err
	}
	s.Logger.Info("successfully added feed to user", "url", url, "name", user.Name)
	urlFeed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		s.Logger.Error("failed to get feed by url", "url", url, "err", err)
		return errors.New("failed to get feed by url")
	}
	var createFeedParams database.CreateFeedFollowParams
	createFeedParams.ID = uuid.New()
	createFeedParams.CreatedAt = curTime
	createFeedParams.UpdatedAt = curTime
	createFeedParams.UserID = user.ID
	createFeedParams.FeedID = urlFeed.ID
	_, err = s.Db.CreateFeedFollow(context.Background(), createFeedParams)
	if err != nil {
		s.Logger.Error("failed to create feed follow", "url", url, "err", err)
		return err
	}
	s.Logger.Info("successfully added feed to user", "url", url, "name", user.Name)
	return nil
}

func handlerFeeds(s *State, _ Command) error {
	allFeeds, err := s.Db.GetAllFeeds(context.Background())
	if err != nil {
		s.Logger.Error("failed to get all feeds", "err", err)
		return err
	}
	for _, feed := range allFeeds {
		fmt.Printf("* name: %v url: %v added_by: %v\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("the add feed handler expects at least one argument the url")
	}
	url := cmd.Args[0]
	urlFeed, err := s.Db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		s.Logger.Error("failed to get feed by url", "url", url, "err", err)
		return errors.New("failed to get feed by url")
	}
	var createFeedParams database.CreateFeedFollowParams
	curTime := time.Now()
	createFeedParams.ID = uuid.New()
	createFeedParams.CreatedAt = curTime
	createFeedParams.UpdatedAt = curTime
	createFeedParams.UserID = user.ID
	createFeedParams.FeedID = urlFeed.ID
	res, err := s.Db.CreateFeedFollow(context.Background(), createFeedParams)
	if err != nil {
		s.Logger.Error("failed to create feed follow", "url", url, "err", err)
		return err
	}
	s.Logger.Info("successfully added feed to user", "url", url, "name", user.Name)
	fmt.Printf("The name of the feed: %v, the username: %v", res[0].FeedName, res[0].UserName)
	return nil
}

func handlerFollowing(s *State, _ Command, user database.User) error {
	var res []database.GetFeedFollowsForUserRow
	res, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		s.Logger.Error("failed to get feed follows", "err", err)
		return err
	}
	if len(res) == 0 {
		s.Logger.Info("no feed follows found for user", "user", user.Name)
		return nil
	}
	fmt.Printf("for the user: %v\n", res[0].Username)
	for _, user := range res {
		fmt.Printf("- %v\n", user.FeedName)
	}
	return nil
}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	if len(cmd.Args) < 1 {
		s.Logger.Error("the add feed handler expects at least one argument the url", "cmd", cmd.Name)
		return errors.New("the unfollow handler expects at least one argument the url")
	}
	feed, err := s.Db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		s.Logger.Error("failed to get feed by url", "url", cmd.Args[0], "err", err)
		return errors.New("failed to get feed by url")
	}
	var deleteUserParams database.DeleteUserAndFeedParams
	deleteUserParams.UserID = user.ID
	deleteUserParams.FeedID = feed.ID
	err = s.Db.DeleteUserAndFeed(context.Background(), deleteUserParams)
	if err != nil {
		s.Logger.Error("failed to delete user and feed", "err", err)
		return err
	}
	s.Logger.Info("successfully unfollowed feed for user", "url", feed.Url, "name", user.Name)
	return nil
}
