package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	err := write(*c)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfgFile, err := os.Open(cfgPath)
	if err != nil {
		return Config{}, nil
	}
	defer cfgFile.Close()

	cfg := Config{}
	decoder := json.NewDecoder(cfgFile)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFileName), nil
}

func write(cfg Config) error {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	cfgFile, err := os.OpenFile(cfgPath, os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	defer cfgFile.Close()

	encoder := json.NewEncoder(cfgFile)
	err = encoder.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
