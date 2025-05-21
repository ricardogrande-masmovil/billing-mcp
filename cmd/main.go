package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	mcpServerSdk "github.com/mark3labs/mcp-go/server"
	"github.com/ricardogrande-masmovil/billing-mcp/api/mcp"
	"github.com/ricardogrande-masmovil/billing-mcp/cmd/di"
	"github.com/ricardogrande-masmovil/billing-mcp/config"
	"github.com/rs/zerolog"
)

const (
	defaultConfigFile  = ".config.yaml"
	migrationFilesPath = "file://database/migrations/schema"
	seedFilesPath      = "file://database/migrations/seeds"
)

var (
	configFile = defaultConfigFile // Use constant for default
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

	// Run seed data if enabled
	if os.Getenv("RUN_SEEDS") == "true" {
		if err := RunSeeds(app.Config, logger); err != nil {
			logger.Error().Err(err).Msg("Failed to run seed data")
			// Decide if you want to exit on seed failure or just log an error
		}
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
	migrateDBURL := cfg.GetMigrateDSN()

	logger.Info().Str("migrationPath", migrationFilesPath).Str("dbURL", migrateDBURL).Msg("Attempting to run migrations")

	m, err := migrate.New(migrationFilesPath, migrateDBURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logger.Error().Err(sourceErr).Msg("Error closing migration source")
		}
		if dbErr != nil {
			logger.Error().Err(dbErr).Msg("Error closing migration database connection")
		}
	}()

	upErr := m.Up()

	if upErr != nil && upErr != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", upErr)
	}

	version, dirty, versionErr := m.Version()

	if versionErr != nil {
		logger.Error().Err(versionErr).Msg("Failed to retrieve migration version status after process.")
	} else {
		if upErr == migrate.ErrNoChange {
			logger.Info().Uint32("version", uint32(version)).Bool("dirty", dirty).Msg("No new migrations to apply. Database is up to date.")
		} else {
			logger.Info().Uint32("version", uint32(version)).Bool("dirty", dirty).Msg("Migrations successfully processed. Database is up to date.")
		}
	}

	return nil
}

func RunSeeds(cfg *config.Config, logger zerolog.Logger) error {
	seedDBURL := cfg.GetMigrateDSN("x-migrations-table=seed_migrations")

	logger.Info().Str("seedPath", seedFilesPath).Str("dbURL", seedDBURL).Msg("Attempting to run seed data with custom table 'seed_migrations'")

	m, err := migrate.New(seedFilesPath, seedDBURL)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance for seeds: %w", err)
	}
	defer func() {
		sourceErr, dbErr := m.Close()
		if sourceErr != nil {
			logger.Error().Err(sourceErr).Msg("Error closing seed source")
		}
		if dbErr != nil {
			logger.Error().Err(dbErr).Msg("Error closing seed database connection")
		}
	}()

	upErr := m.Up()

	if upErr != nil && upErr != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply seeds: %w", upErr)
	}

	if upErr == migrate.ErrNoChange {
		logger.Info().Msg("No new seeds to apply.")
	} else {
		logger.Info().Msg("Seed data successfully applied.")
	}

	version, dirty, versionErr := m.Version()
	if versionErr != nil {
		logger.Error().Err(versionErr).Msg("Failed to retrieve seed version status after process.")
	} else {
		logger.Info().Uint32("version", uint32(version)).Bool("dirty", dirty).Msg("Seed data successfully processed.")
	}

	return nil
}