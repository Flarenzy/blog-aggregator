package command

import (
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"github.com/Flarenzy/blog-aggregator/internal/database"
	"github.com/Flarenzy/blog-aggregator/internal/rss"
	"github.com/google/uuid"
	"html"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))
	return cmds
}

func (c *Commands) register(name string, f func(*State, Command) error) {
	_, ok := c.registered[name]
	if !ok {
		c.registered[name] = f
		slog.Debug("registering command", "name", name)
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

func fetchFeed(ctx context.Context, feedUrl string) (*rss.RSSFeed, error) {
	withContext, err := http.NewRequestWithContext(ctx, "GET", feedUrl, nil)
	if err != nil {
		return nil, err
	}
	withContext.Header.Set("User-Agent", "gator")
	resp, err := http.DefaultClient.Do(withContext)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var feed *rss.RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return feed, nil
}

func unescapeXML(rss *rss.RSSFeed) (*rss.RSSFeed, error) {
	if rss == nil {
		return nil, errors.New("nil rss feed")
	}
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, entry := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(entry.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(entry.Description)
	}
	return rss, nil
}

func scrapeFeeds(s *State) {
	s.Logger.Debug("scraping feeds")
	user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		s.Logger.Error("Error getting user", "error", err)
		return
	}

	nextFeed, err := s.Db.GetNextFeedToFetch(context.Background(), user.ID)
	if err != nil {
		s.Logger.Info("No more feeds to fetch", "error", err)
		return
	}
	feed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		s.Logger.Info("Error fetching feed", "error", err, "url", nextFeed.Url)
		return
	}
	if feed == nil {
		s.Logger.Info("nil feed", "url", nextFeed.Url)
		return
	}
	var markFeedFetched database.MarkFeedFetchedParams
	markFeedFetched.UpdatedAt = time.Now()
	markFeedFetched.LastFetchedAt = sql.NullTime{
		Time:  markFeedFetched.UpdatedAt,
		Valid: true,
	}
	markFeedFetched.ID = nextFeed.ID

	err = s.Db.MarkFeedFetched(context.Background(), markFeedFetched)
	if err != nil {
		s.Logger.Error("Error marking feed as fetched", "error", err, "id", nextFeed.ID)
		return
	}
	unescapedFeed, err := unescapeXML(feed)
	if err != nil {
		s.Logger.Error("Error unescaping feed", "error", err, "id", nextFeed.ID)
		return
	}
	if unescapedFeed == nil {
		s.Logger.Error("Error got nil pointer unescapedFeed", "error", err, "id", nextFeed.ID)
		return
	}
	for _, item := range unescapedFeed.Channel.Item {
		if item.Title == "" {
			continue
		}
		var createPostParams database.CreatePostParams
		createPostParams.ID = uuid.New()
		createPostParams.CreatedAt = time.Now()
		createPostParams.UpdatedAt = createPostParams.CreatedAt
		parsedDate, err := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", item.PubDate)
		if err != nil {
			s.Logger.Error("Error parsing date format", "error", err, "date", item.PubDate)
			continue
		}
		createPostParams.PublishedAt = parsedDate
		createPostParams.FeedID = nextFeed.ID
		createPostParams.Url = item.Link
		createPostParams.Title = item.Title
		var nulDescription sql.NullString
		if item.Description != "" {
			nulDescription.Valid = true
			nulDescription.String = item.Description
		} else {
			nulDescription.Valid = false
			nulDescription.String = ""
		}
		createPostParams.Description = nulDescription
		_, err = s.Db.CreatePost(context.Background(), createPostParams)
		if err != nil {
			s.Logger.Error("Error creating post", "error", err, "id", nextFeed.ID, "url", item.Link)
			continue
		}
		s.Logger.Info("Created post", "Title", item.Title, "url", item.Link)
	}

}
