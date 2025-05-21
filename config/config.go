package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DatabaseConfig struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"dbname"`
	SSLMode    string `yaml:"sslmode"`
	MaxRetries int    `yaml:"maxRetries"`
}

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	LogLevel string         `yaml:"logLevel"`
	Version  string         `yaml:"version"`
	RunSeeds bool           `yaml:"runSeeds"` // Added RunSeeds flag
}

// LoadConfig loads configuration from the given YAML file path.
func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config YAML: %w", err)
	}

	// Set default values
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080" // Default port
	}
	if cfg.Database.SSLMode == "" {
		cfg.Database.SSLMode = "disable" // Default SSLMode
	}
	if cfg.Database.MaxRetries == 0 {
		cfg.Database.MaxRetries = 3 // Default MaxRetries
	}
	
	return &cfg, nil
}

// GetDSN constructs the Data Source Name (DSN) for connecting to the PostgreSQL database.
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetMigrateDSN constructs the Data Source Name (DSN) for golang-migrate, which needs to be in a URL format.
//
//	This is different from the DSN used by gorm. Example: postgresql://user:password@host:port/dbname?sslmode=disable
func (c *Config) GetMigrateDSN(params ...string) string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
	if len(params) > 0 {
		dsn += "&" + params[0]
	}
	return dsn
}
