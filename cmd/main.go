package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	mcpServerSdk "github.com/mark3labs/mcp-go/server"
	"github.com/ricardogrande-masmovil/billing-mcp/api/mcp"
	"github.com/ricardogrande-masmovil/billing-mcp/cmd/di"
	"github.com/ricardogrande-masmovil/billing-mcp/config"
	"github.com/rs/zerolog"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	configFile = ".config.yaml" // Default
)

func main() {
	if cp := os.Getenv("CONFIG_PATH"); cp != "" {
		configFile = cp
	}

	app, cleanup, err := di.InitializeApp(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	logger := app.Logger
	logger.Info().Msg("Successfully initialized application dependencies")

	// Run database migrations
	if err := RunMigrations(app.Config, logger); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	logger.Info().Msg("Starting the application...")

	ctx := context.Background()

	exitChannel := make(chan bool, 1)

	go InitMCP(ctx, app.Echo, app.MCPServer, app.MCPServerAPI, app.Config, app.Logger, exitChannel)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	logger.Info().Msg("Received shutdown signal, shutting down...")
	exitChannel <- true
	<-exitChannel
	logger.Info().Msg("Application shutdown complete.")
}

func InitMCP(ctx context.Context, e *echo.Echo, sdkServer *mcpServerSdk.MCPServer, appMCPServer *mcp.MCPServer, cfg *config.Config, logger zerolog.Logger, exitChan chan bool) {
	err := mcp.Setup(e, sdkServer, appMCPServer)
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to setup MCP server")
	}

	serverAddr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	logger.Info().Str("address", serverAddr).Msg("Starting MCP server...")

	go func() {
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed { // Check for http.ErrServerClosed
			logger.Error().Err(err).Msg("MCP server failed to start")
		} else if err == http.ErrServerClosed {
			logger.Info().Msg("MCP server stopped gracefully.")
		} else {
			logger.Info().Msg("MCP server stopped.")
		}
	}()

	<-exitChan

	logger.Warn().Msg("Shutting down MCP server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Failed to shutdown MCP server gracefully")
	} else {
		logger.Info().Msg("MCP server shutdown complete")
	}
	exitChan <- true
}

func RunMigrations(cfg *config.Config, logger zerolog.Logger) error {
	migrationPath := "file://database/migrations"

	// For postgres, we need to ensure the DSN is in the format expected by the migrate tool
	// which is slightly different from gorm's DSN. Specifically, it needs to be a URL.
	// Example: postgresql://user:password@host:port/dbname?sslmode=disable
	migrateDBURL := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode)

	logger.Info().Str("migrationPath", migrationPath).Str("dbURL", migrateDBURL).Msg("Attempting to run migrations")

	m, err := migrate.New(migrationPath, migrateDBURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get migration version after applying")
	} else {
		logger.Info().Uint32("version", uint32(version)).Bool("dirty", dirty).Msg("Migrations applied successfully")
	}

	sourceErr, dbErr := m.Close()
	if sourceErr != nil {
		logger.Error().Err(sourceErr).Msg("Error closing migration source")
	}
	if dbErr != nil {
		logger.Error().Err(dbErr).Msg("Error closing migration database connection")
	}

	return nil
}
