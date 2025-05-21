package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig_Success(t *testing.T) {
	configPath := "../.config.example.yaml" // Relative to the config package directory

	expectedConfig := &Config{
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Database: DatabaseConfig{
			Host:       "localhost",
			Port:       5432,
			User:       "billing_user",
			Password:   "yoursecurepassword",
			DBName:     "billing_db",
			SSLMode:    "disable",
			MaxRetries: 3, // Added MaxRetries
		},
		LogLevel: "info",
		Version:  "0.0.1",
		RunSeeds: false, // Assuming default is false and not set in .config.example.yaml
	}

	cfg, err := LoadConfig(configPath)

	require.NoError(t, err, "LoadConfig() should not return an error")
	assert.Equal(t, expectedConfig, cfg, "Loaded config should match expected config")
}

func TestLoadConfig_FileNotExist(t *testing.T) {
	_, err := LoadConfig("non_existent_config.yaml")
	assert.Error(t, err, "LoadConfig() should return an error for a non-existent file")
}

func TestLoadConfig_Defaults(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "temp_config.yaml")

	// Create a temporary config file with missing server port and sslmode
	content := []byte(`
server:
  host: "testhost"
database:
  host: "dbhost"
  port: 1234
  user: "testuser"
  password: "testpass"
  dbname: "testdb"
logLevel: "debug"
version: "0.0.1"
`)
	require.NoError(t, os.WriteFile(tempFile, content, 0600), "Failed to write temp config file")

	cfg, err := LoadConfig(tempFile)
	require.NoError(t, err, "LoadConfig() should not return an error for valid temp file")

	assert.Equal(t, "8080", cfg.Server.Port, "Default server port should be applied")
	assert.Equal(t, "disable", cfg.Database.SSLMode, "Default SSL mode should be applied")
	assert.Equal(t, 3, cfg.Database.MaxRetries, "Default MaxRetries should be applied")
	assert.False(t, cfg.RunSeeds, "Default RunSeeds should be false")

	// Check other values are loaded correctly
	assert.Equal(t, "testhost", cfg.Server.Host)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, 1234, cfg.Database.Port)
}

func TestLoadConfig_RunSeedsTrue(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "temp_config_runseeds.yaml")

	content := []byte(`
server:
  host: "testhost"
database:
  host: "dbhost"
  port: 1234
  user: "testuser"
  password: "testpass"
  dbname: "testdb"
logLevel: "debug"
version: "0.0.1"
runSeeds: true
`)
	require.NoError(t, os.WriteFile(tempFile, content, 0600), "Failed to write temp config file")

	cfg, err := LoadConfig(tempFile)
	require.NoError(t, err, "LoadConfig() should not return an error for valid temp file")

	assert.True(t, cfg.RunSeeds, "RunSeeds should be true when set in config")
	assert.Equal(t, "8080", cfg.Server.Port, "Default server port should be applied")
	assert.Equal(t, "disable", cfg.Database.SSLMode, "Default SSL mode should be applied")
	assert.Equal(t, 3, cfg.Database.MaxRetries, "Default MaxRetries should be applied")
}
