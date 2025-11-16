package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

func NewConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	configContent := os.ExpandEnv(string(data))

	var cfg Config
	if err := yaml.Unmarshal([]byte(configContent), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	cfg.setDefaults()

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) setDefaults() {
	if c.Server.Port == "" {
		c.Server.Port = "8080"
	}
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Database.Port == "" {
		c.Database.Port = "5432"
	}
	if c.Database.SSLMode == "" {
		c.Database.SSLMode = "disable"
	}
}

func (c *Config) validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

func (d *DatabaseConfig) GetDSN() string {
	if d.Password == "" {
		return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
			d.User, d.Host, d.Port, d.DBName, d.SSLMode)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}
