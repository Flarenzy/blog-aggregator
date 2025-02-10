package main

import (
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"log/slog"
)

func main() {
	c, err := config.Read("")
	if err != nil {
		slog.Error("error reading config", "err", err)
	}
	err = c.SetUser("Flarenzy")
	if err != nil {
		slog.Error("error setting user", "err", err)
	}
	fmt.Printf("DB url: %v\nUser: %v\n", c.DbUrl, c.CurrentUserName)

}
