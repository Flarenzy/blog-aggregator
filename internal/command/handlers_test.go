package command

import (
	"github.com/Flarenzy/blog-aggregator/internal/config"
	"os"
	"path/filepath"
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
	s := NewState(c)
	err := handlerLogin(s, Command{
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
	s := NewState(c)
	err := handlerLogin(s,
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
		s := NewState(c)
		args := []string{ca}
		err := handlerLogin(s,
			Command{
				"login",
				args,
			})
		if err != nil && s.Config.CurrentUserName != ca {
			t.Errorf("expected no error got: %v", err.Error())
			t.Fatal()
		}
	}
}
