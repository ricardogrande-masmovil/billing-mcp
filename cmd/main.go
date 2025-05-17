package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ricardogrande-masmovil/billing-mcp/api/mcp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mcpServerSdk "github.com/mark3labs/mcp-go/server"
)

var logger = log.With().Str("module", "main").Logger()

func main() {
	ctx := context.Background()

	// Initialize the logger
	InitLogger("debug")

	logger.Info().Msg("Starting the application...")

	exitChannel := make(chan bool, 1)

	go InitMCP(ctx, exitChannel)

	// Run until SIGTERM or SIGINTERRUPT is received
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	<-quit

	logger.Info().Msg("Received shutdown signal, shutting down...")
	exitChannel <- true
}

func InitMCP(ctx context.Context, exitChan chan bool) {
	e := echo.New()
	s := mcpServerSdk.NewMCPServer("billing", "0.0.1")

	err := mcp.Setup(e, s)
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to setup MCP server")
	}

	go func() {
		// TODO: define port in config
		err := e.Start(":8080")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to start MCP server")
		}
	}()

	// Shutdown MCP server gracefully and propagate the exit signal
	exit := <-exitChan
	exitChan <- exit

	logger.Warn().Msg("Shutting down MCP server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to shutdown MCP server")
	}
	logger.Info().Msg("MCP server shutdown complete")
}

func InitLogger(level string) {
	// TODO: setup logging level in config
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse log level")
	}
	zerolog.SetGlobalLevel(logLevel)
}
