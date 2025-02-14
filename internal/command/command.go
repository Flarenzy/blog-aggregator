package command

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"github.com/Flarenzy/blog-aggregator/internal/database"
	"github.com/Flarenzy/blog-aggregator/internal/rss"
	"html"
	"io"
	"log/slog"
	"net/http"
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
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
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
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for i, entry := range rss.Channel.Item {
		rss.Channel.Item[i].Title = html.UnescapeString(entry.Title)
		rss.Channel.Item[i].Description = html.UnescapeString(entry.Description)
	}
	return rss, nil
}
