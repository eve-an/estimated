package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	SessionChannelSize int    `json:"session_channel_size"`
	ServerAddress      string `json:"server_address"`
}

func (c *Config) validate() error {
	var errs []error

	if c.ServerAddress == "" {
		errs = append(errs, errors.New("field 'server_address' must not be empty"))
	}

	return errors.Join(errs...)
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

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
