package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	SessionChannelSize int `json:"session_channel_size"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decoding config file: %w", err)
	}

	return &cfg, nil
}
