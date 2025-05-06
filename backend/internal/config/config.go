package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

type Config struct {
	SessionChannelSize int      `json:"session_channel_size"`
	ServerAddress      string   `json:"server_address"`
	ServerTimeout      Duration `json:"server_timeout"`
}

func (c *Config) validate() error {
	var errs []error

	if c.ServerAddress == "" {
		errs = append(errs, errors.New("field 'server_address' must not be empty"))
	}

	if c.ServerTimeout == 0 {
		errs = append(errs, errors.New("field 'server_timeout' must not be empty"))
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
