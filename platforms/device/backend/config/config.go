package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the device platform configuration
type Config struct {
	Mode     string `yaml:"mode"`     // "gateway" or "hmi"
	Runtime  struct {
		Storage struct {
			Shards    int    `yaml:"shards"`
			ShardSize int    `yaml:"shard_size"`
		} `yaml:"storage"`
	} `yaml:"runtime"`
	WebSocket struct {
		Enabled         bool   `yaml:"enabled"`
		Port            int    `yaml:"port"`
		Path            string `yaml:"path"`
		ReadBufferSize  int    `yaml:"read_buffer_size"`
		WriteBufferSize int    `yaml:"write_buffer_size"`
		PingPeriod      int    `yaml:"ping_period"`
	} `yaml:"websocket"`
	HTTP struct {
		Enabled      bool   `yaml:"enabled"`
		Port         int    `yaml:"port"`
		ReadTimeout  int    `yaml:"read_timeout"`
		WriteTimeout int    `yaml:"write_timeout"`
	} `yaml:"http"`
}

// Load loads configuration from YAML file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// validate validates the configuration
func validate(cfg *Config) error {
	if cfg.Mode != "gateway" && cfg.Mode != "hmi" {
		return fmt.Errorf("invalid mode: %s (must be 'gateway' or 'hmi')", cfg.Mode)
	}

	if cfg.WebSocket.Enabled && cfg.WebSocket.Port <= 0 {
		return fmt.Errorf("invalid WebSocket port: %d", cfg.WebSocket.Port)
	}

	if cfg.HTTP.Enabled && cfg.HTTP.Port <= 0 {
		return fmt.Errorf("invalid HTTP port: %d", cfg.HTTP.Port)
	}

	return nil
}
