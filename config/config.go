package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	TelegramToken string `json:"telegram_token"`
	DatabaseURL   string `json:"database_url"`
}

func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
