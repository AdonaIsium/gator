package internal

import (
	"encoding/json"
	"os"
)

const configFileName = "~/.gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (*Config, error) {
	data, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) SetUser(username string) error {
	cfg.CurrentUserName = username

	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	os.WriteFile(configFileName, data, 0644)

	return nil
}

func getConfigFilePath() (string, error) {
	return configFileName, nil
}
