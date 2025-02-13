package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/command"
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"github.com/Flarenzy/blog-aggregator/internal/database"
	"github.com/Flarenzy/blog-aggregator/internal/logging"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	// Gracefully handle interrupt signals

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		os.Exit(0)
	}()
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	c, err := config.Read("")
	if err != nil {
		slog.Error("error reading config", "err", err)
	}
	logger, logFile, err := logging.NewLogger("gator.log", slog.LevelDebug)
	if err != nil {
		slog.Error("error creating logger", "err", err)
		os.Exit(1)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			slog.Error("error closing log file", "err", err)
			os.Exit(1)
		}
	}(logFile)
	db, err := sql.Open("postgres", c.DbUrl)
	dbQueries := database.New(db)
	s := command.NewState(&c, dbQueries, logger)
	cmds := command.NewCommands()
	args := os.Args[1:]
	fmt.Println(args)
	if len(args) < 1 {
		fmt.Println("we need at least two arguments")
		os.Exit(1)
	}
	cmd := command.Command{
		Name: args[0],
		Args: args[1:],
	}
	err = cmds.Run(s, cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
