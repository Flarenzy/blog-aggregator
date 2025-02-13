package command

import (
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"github.com/Flarenzy/blog-aggregator/internal/logging"
	"log/slog"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func tempConfig(t *testing.T) *config.Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("error getting user home dir %v", err)
		t.Fatal()
	}
	f, err := os.CreateTemp(homeDir, ".gatorconfig.json")
	if err != nil {
		t.Errorf("Failed creating temp file %v", err)
		t.Fatal()
	}
	_, err = f.Write([]byte("{\"db_url\":\"postgres://example\"}"))
	if err != nil {
		t.Errorf("Failed writing to temp file %v", err)
		t.Fatal()
	}
	baseName := filepath.Base(f.Name())
	c, err := config.Read(baseName)
	if err != nil {
		t.Errorf("Failed reading config %v", err)
		t.Fatal()
	}
	return &c
}

func TestLoginNotEnoughArgs(t *testing.T) {
	t.Parallel()
	c := tempConfig(t)
	logger, f, err := logging.NewLogger("gator-temp1.log", slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(f)
	defer func() {
		err := os.Remove("gator-temp1.log")
		if err != nil {
			t.Fatal(err)
		}
	}()
	s := NewState(c, logger)
	err = handlerLogin(s, Command{
		"login",
		[]string{},
	})
	if err != nil && err.Error() != "the login handler expects a single argument, the username" {
		t.Errorf("expected a different error got: %v", err.Error())
		t.Fatal()
	}
}

func TestLoginEmptyUsername(t *testing.T) {
	t.Parallel()
	c := tempConfig(t)
	logger, f, err := logging.NewLogger("gator-temp2.log", slog.LevelDebug)
	if err != nil {
		t.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(f)
	defer func() {
		err := os.Remove("gator-temp2.log")
		if err != nil {
			t.Fatal(err)
		}
	}()
	s := NewState(c, logger)
	err = handlerLogin(s,
		Command{
			"login",
			[]string{""},
		})
	if err != nil && err.Error() != "user is empty" {
		t.Errorf("expected a different error got: %v", err.Error())
		t.Fatal()
	}
}

func TestLoginInvalidUsername(t *testing.T) {
	t.Parallel()
	cases := []string{"Bob", "Frank", "John"}
	for _, ca := range cases {
		c := tempConfig(t)
		randNum := strconv.Itoa(rand.Int())
		logger, f, err := logging.NewLogger("gator-temp"+randNum+".log", slog.LevelDebug)
		if err != nil {
			t.Fatal(err)
		}

		s := NewState(c, logger)
		args := []string{ca}
		err = handlerLogin(s,
			Command{
				"login",
				args,
			})
		if err != nil && s.Config.CurrentUserName != ca {
			t.Errorf("expected no error got: %v", err.Error())
			f.Close()
			os.Remove("gator-temp.log")
			t.Fatal()
		}
		f.Close()
		os.Remove(f.Name())
	}
}
