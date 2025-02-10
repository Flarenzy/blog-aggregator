package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
	fp              string
}

func newConfig() *Config {
	return &Config{}
}

func Read(fileName string) (Config, error) {
	c := newConfig()
	path, err := getConfigFilepath(fileName)
	c.fp = path
	if err != nil {
		return Config{}, err
	}
	file, err := os.ReadFile(path)

	if err != nil && err != io.EOF {
		return Config{}, err
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

func (c *Config) SetUser(user string) error {
	if user == "" {
		return errors.New("user is empty")
	}
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		return err
	}

	return nil
}

func write(c Config) error {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return err
	}
	var path string
	if c.fp != "" {
		path = c.fp
	} else {
		path, err = getConfigFilepath("")
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilepath(fileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	if fileName == "" {
		fileName = filepath.Join(home, configFileName)
	} else {
		fileName = filepath.Join(home, fileName)
	}
	return fileName, nil
}
