package main

import (
	"context"
	"fmt"
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
)

var (
	configFile = ".config.yaml"
)

func main() {
	app, cleanup, err := di.InitializeApp(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	logger := app.Logger
	logger.Info().Msg("Successfully initialized application dependencies")
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
		if err := e.Start(serverAddr); err != nil {
			logger.Error().Err(err).Msg("Failed to start MCP server")
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
