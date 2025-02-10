package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConfig(t *testing.T) {
	c := newConfig()
	if c.DbUrl != "" {
		t.Errorf("db url should be empty")
	}
	if c.CurrentUserName != "" {
		t.Errorf("current user name should be empty")
	}
}

func TestNewConfigFromFile(t *testing.T) {

	cases := []string{
		"{\"db_url\":\"postgres://example\", \"current_user_name\":\"Flarenzy\"}",
		"{\"db_url\":\"postgres://example\"}",
		"{\"current_user_name\":\"Flarenzy\"}",
		"{}",
	}
	expected := [][]string{
		{"postgres://example", "Flarenzy"},
		{"postgres://example", ""},
		{"", "Flarenzy"},
		{"", ""},
	}
	for i, ca := range cases {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Errorf("error getting user home dir %v", err)
			t.Fatal()
		}
		f, err := os.CreateTemp(homeDir, configFileName)
		if err != nil {
			t.Errorf("Failed creating temp file %v", err)
			t.Fatal()
		}

		_, err = f.Write([]byte(ca))
		if err != nil {
			t.Errorf("i: %v, ca: %v", i, ca)
			t.Errorf("Failed writing to temp file %v", err)
			t.Fatal()
		}
		baseName := filepath.Base(f.Name())
		c, err := Read(baseName)
		if err != nil {
			t.Errorf("i: %v, ca: %v", i, ca)
			t.Errorf("Failed reading temp file %v", err)
			t.Fatal()
		}
		if c.DbUrl != expected[i][0] {
			t.Errorf("db url should be postgres://example")

		}
		if c.CurrentUserName != expected[i][1] {
			t.Errorf("current user name should be Flarenzy")
		}
		err = os.Remove(f.Name())
		if err != nil {
			t.Errorf("Failed removing temp file %v", err)
			t.Fatal()
		}
		err = f.Close()
		if err != nil {
			t.Errorf("Failed closing temp file %v", err)
			t.Fatal()
		}
	}
}

func TestGetConfigFilePath(t *testing.T) {
	cases := []string{"", "random.json"}
	expected := []string{configFileName, "random.json"}
	for i, ca := range cases {
		fp, err := getConfigFilepath(ca)
		if err != nil {
			t.Errorf("Failed getting config file path %v", err)
			t.Fatal()
		}
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Errorf("Failed getting user home dir %v", err)
			t.Fatal()
		}
		expectedPath := filepath.Join(homeDir, expected[i])
		if fp != expectedPath {
			t.Errorf("config file path does not match expected path")
			t.Fatal()
		}
	}
}

func TestSetUser(t *testing.T) {
	cases := []string{"Flarenzy", "John", "Boots"}
	expected := []string{"Flarenzy", "John", "Boots"}

	for i, ca := range cases {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			t.Errorf("error getting user home dir %v", err)
			t.Fatal()
		}
		f, err := os.CreateTemp(homeDir, configFileName)
		if err != nil {
			t.Errorf("Failed creating temp file %v", err)
			t.Fatal()
		}

		_, err = f.Write([]byte("{\"db_url\":\"postgres://example\"}"))
		if err != nil {
			t.Errorf("i: %v, ca: %v", i, ca)
			t.Errorf("Failed writing to temp file %v", err)
			t.Fatal()
		}
		baseName := filepath.Base(f.Name())
		c, err := Read(baseName)
		if err != nil {
			t.Errorf("i: %v, ca: %v", i, ca)
			t.Errorf("Failed reading temp file %v", err)
			t.Fatal()
		}
		err = c.SetUser(ca)
		if err != nil {
			t.Errorf("i: %v, ca: %v", i, ca)
			t.Errorf("Failed setting user %v", err)
			t.Fatal()
		}
		if c.CurrentUserName != expected[i] {
			t.Errorf("current user name should be %v, but %v", expected[i], c.CurrentUserName)
		}
		err = os.Remove(f.Name())
		if err != nil {
			t.Errorf("Failed removing temp file %v", err)
			t.Fatal()
		}
		err = f.Close()
		if err != nil {
			t.Errorf("Failed closing temp file %v", err)
			t.Fatal()
		}
	}
}
