// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ricardogrande-masmovil/billing-mcp/api"
	"github.com/ricardogrande-masmovil/billing-mcp/api/mcp"
	"github.com/ricardogrande-masmovil/billing-mcp/config"
	domain2 "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain"
	persistence2 "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence/sql"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/ports"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain"
	persistence3 "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/infrastructure/persistence"
	sql2 "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/infrastructure/persistence/sql"
	ports2 "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/ports"
	"github.com/ricardogrande-masmovil/billing-mcp/pkg/persistence"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// Injectors from wire.go:

func InitializeApp(configFile string) (*App, func(), error) {
	config, err := ProvideConfig(configFile)
	if err != nil {
		return nil, nil, err
	}
	logger := ProvideLogger(config)
	db, cleanup, err := ProvideDB(config, logger)
	if err != nil {
		return nil, nil, err
	}
	echo := ProvideEcho()
	mcpServer := ProvideMCP(config)
	healthController := ProvideHealthController()
	invoiceSqlClient := ProvideInvoiceSqlClient(db, config)
	invoiceSqlConverter := ProvideInvoiceSqlConverter()
	repository := ProvideInvoicePersistenceRepository(invoiceSqlClient, invoiceSqlConverter)
	service := ProvideInvoiceDomainService(repository)
	invoicesController := ProvideInvoicesController(service)
	movementSqlClient := ProvideMovementSqlClient(db, logger)
	movementConverter := ProvideMovementConverter()
	movementRepository := ProvideMovementRepository(movementSqlClient, movementConverter, logger)
	movementService := ProvideMovementService(logger, movementRepository)
	movementsController := ProvideMovementsController(movementService, logger)
	mcpMCPServer := ProvideMCPServerAPI(healthController, invoicesController, movementsController)
	app := &App{
		Config:              config,
		Logger:              logger,
		DB:                  db,
		Echo:                echo,
		MCPServer:           mcpServer,
		MCPServerAPI:        mcpMCPServer,
		HealthController:    healthController,
		InvoicesController:  invoicesController,
		MovementsController: movementsController,
		MovementsService:    movementService,
	}
	return app, func() {
		cleanup()
	}, nil
}

// wire.go:

// App holds the application's dependencies.
type App struct {
	Config              *config.Config
	Logger              zerolog.Logger
	DB                  *gorm.DB
	Echo                *echo.Echo
	MCPServer           *server.MCPServer
	MCPServerAPI        *mcp.MCPServer // Added field for the API specific MCP server
	HealthController    mcp.HealthController
	InvoicesController  mcp.InvoicesController
	MovementsController mcp.MovementsController
	MovementsService    domain.MovementService
}

// --- Core Providers ---
func ProvideConfig(filePath string) (*config.Config, error) {
	return config.LoadConfig(filePath)
}

func ProvideLogger(cfg *config.Config) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse log level for DI")
	}
	return log.Logger.Level(logLevel)
}

func ProvideDB(cfg *config.Config, logger zerolog.Logger) (*gorm.DB, func(), error) {
	dsn := cfg.GetDSN()
	db, err := persistence.NewSqlClient(dsn)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		sqlDB, dbErr := db.DB()
		if dbErr != nil {
			logger.Error().Err(dbErr).Msg("Failed to get underlying sql.DB for cleanup")
			return
		}
		if err := sqlDB.Close(); err != nil {
			logger.Error().Err(err).Msg("Failed to close database connection")
		} else {
			logger.Info().Msg("Database connection closed successfully.")
		}
	}
	return db, cleanup, nil
}

func ProvideEcho() *echo.Echo {
	return echo.New()
}

func ProvideMCP(cfg *config.Config) *server.MCPServer {
	return server.NewMCPServer("billing-mcp", cfg.Version)
}

// Provider for the API specific MCPServer
func ProvideMCPServerAPI(healthController mcp.HealthController, invoicesController mcp.InvoicesController, movementsController mcp.MovementsController) *mcp.MCPServer {
	return mcp.NewMCPServer(healthController, invoicesController, movementsController)
}

func ProvideHealthController() mcp.HealthController {
	return api.NewHealthController()
}

// --- Invoice Feature Providers ---
func ProvideInvoiceSqlClient(db *gorm.DB, cfg *config.Config) sql.InvoiceSqlClient {
	return sql.NewInvoiceSqlClient(db, cfg.Database.MaxRetries)
}

func ProvideInvoiceSqlConverter() sql.InvoiceSqlConverter {
	return sql.NewInvoiceSqlConverter()
}

func ProvideInvoicePersistenceRepository(client sql.InvoiceSqlClient, converter sql.InvoiceSqlConverter) persistence2.Repository {
	return persistence2.NewRepository(client, converter)
}

func ProvideInvoiceDomainService(repo domain2.Repository) domain2.Service {
	return domain2.NewService(repo)
}

func ProvideInvoicePortsService(domainService domain2.Service) ports.InvoiceService {
	return domainService
}

func ProvideInvoicesController(service ports.InvoiceService) mcp.InvoicesController {
	return ports.NewController(service)
}

// --- Movement Feature Providers ---
func ProvideMovementsController(movementService domain.MovementService, logger zerolog.Logger) mcp.MovementsController {
	return ports2.NewMCPMovementsHandler(movementService, logger)
}

func ProvideMovementSqlClient(db *gorm.DB, logger zerolog.Logger) *sql2.MovementSqlClient {
	return sql2.NewMovementSqlClient(db, logger)
}

func ProvideMovementConverter() *sql2.MovementConverter {
	return sql2.NewMovementConverter()
}

func ProvideMovementRepository(client *sql2.MovementSqlClient, converter *sql2.MovementConverter, logger zerolog.Logger) domain.MovementRepository {
	return persistence3.NewMovementSQLRepository(client, converter, logger)
}

func ProvideMovementService(logger zerolog.Logger, repo domain.MovementRepository) domain.MovementService {
	return *domain.NewMovementService(logger, repo)
}

// --- Provider Sets ---
var CoreSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDB,
	ProvideEcho,
	ProvideMCP,
	ProvideMCPServerAPI,
	ProvideHealthController,
)

var InvoiceFeatureSet = wire.NewSet(
	ProvideInvoiceSqlClient,
	ProvideInvoiceSqlConverter,
	ProvideInvoicePersistenceRepository, wire.Bind(new(domain2.Repository), new(persistence2.Repository)), ProvideInvoiceDomainService, wire.Bind(new(ports.InvoiceService), new(domain2.Service)), ProvideInvoicesController,
)

var MovementFeatureSet = wire.NewSet(
	ProvideMovementSqlClient,
	ProvideMovementConverter,
	ProvideMovementRepository,
	ProvideMovementService,
	ProvideMovementsController,
)

var AppSet = wire.NewSet(
	CoreSet,
	InvoiceFeatureSet,
	MovementFeatureSet, wire.Struct(new(App), "*"),
)
