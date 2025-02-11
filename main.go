package main

import (
	"fmt"
	"github.com/Flarenzy/blog-aggregator/internal/command"
	"github.com/Flarenzy/blog-aggregator/internal/config"
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
	c, err := config.Read("")
	if err != nil {
		slog.Error("error reading config", "err", err)
	}
	s := command.NewState(&c)
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
