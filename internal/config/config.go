package config

import (
	"os";
	"fmt";
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DbUrl string `json:db_url`
	CurrentUserName string `json:current_user_name`
}

func Read() (Config, error) {
	cfgPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	
	return Config{}, nil
}

func SetUser(username string) error {
	return nil
}

func getConfigFilePath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", dir, configFileName), nil
}

func write(cfg Config) error {
	return nil
}


