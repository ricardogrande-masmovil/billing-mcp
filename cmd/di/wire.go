//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	mcpServerSdk "github.com/mark3labs/mcp-go/server"

	"github.com/ricardogrande-masmovil/billing-mcp/api"
	mcpAPI "github.com/ricardogrande-masmovil/billing-mcp/api/mcp"
	"github.com/ricardogrande-masmovil/billing-mcp/config"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain"
	invoicePersistence "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence"
	invoiceSQL "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence/sql"
	invoicePorts "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/ports"
	movementsDomain "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain"
	movementsPersistence "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/infrastructure/persistence"
	pkgPersistence "github.com/ricardogrande-masmovil/billing-mcp/pkg/persistence"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// App holds the application's dependencies.
type App struct {
	Config             *config.Config
	Logger             zerolog.Logger
	DB                 *gorm.DB
	Echo               *echo.Echo
	MCPServer          *mcpServerSdk.MCPServer
	MCPServerAPI       *mcpAPI.MCPServer // Added field for the API specific MCP server
	HealthController   mcpAPI.HealthController
	InvoicesController mcpAPI.InvoicesController
	MovementsService   movementsDomain.MovementService
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
	db, err := pkgPersistence.NewSqlClient(dsn)
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

func ProvideMCP(cfg *config.Config) *mcpServerSdk.MCPServer {
	return mcpServerSdk.NewMCPServer("billing-mcp", cfg.Version)
}

// Provider for the API specific MCPServer
func ProvideMCPServerAPI(healthController mcpAPI.HealthController, invoicesController mcpAPI.InvoicesController) *mcpAPI.MCPServer {
	return mcpAPI.NewMCPServer(healthController, invoicesController)
}

func ProvideHealthController() mcpAPI.HealthController {
	return api.NewHealthController()
}

// --- Invoice Feature Providers ---
func ProvideInvoiceSqlClient(db *gorm.DB, cfg *config.Config) invoiceSQL.InvoiceSqlClient {
	return invoiceSQL.NewInvoiceSqlClient(db, cfg.Database.MaxRetries)
}

func ProvideInvoiceSqlConverter() invoiceSQL.InvoiceSqlConverter {
	return invoiceSQL.NewInvoiceSqlConverter()
}

func ProvideInvoicePersistenceRepository(client invoiceSQL.InvoiceSqlClient, converter invoiceSQL.InvoiceSqlConverter) invoicePersistence.Repository {
	return invoicePersistence.NewRepository(client, converter)
}

func ProvideInvoiceDomainService(repo domain.Repository) domain.Service {
	return domain.NewService(repo)
}

func ProvideInvoicePortsService(domainService domain.Service) invoicePorts.InvoiceService {
	return domainService
}

func ProvideInvoicesController(service invoicePorts.InvoiceService) mcpAPI.InvoicesController {
	return invoicePorts.NewController(service)
}

// --- Movement Feature Providers ---
func ProvideMovementRepository(db *gorm.DB, logger zerolog.Logger) movementsDomain.MovementRepository {
	return movementsPersistence.NewMovementSQLRepository(db, logger)
}

func ProvideMovementService(logger zerolog.Logger, repo movementsDomain.MovementRepository) movementsDomain.MovementService {
	return *movementsDomain.NewMovementService(logger, repo)
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
	ProvideInvoicePersistenceRepository,
	wire.Bind(new(domain.Repository), new(invoicePersistence.Repository)),
	ProvideInvoiceDomainService,
	wire.Bind(new(invoicePorts.InvoiceService), new(domain.Service)),
	ProvideInvoicesController,
)

var MovementFeatureSet = wire.NewSet(
	ProvideMovementRepository,
	ProvideMovementService,
)

var AppSet = wire.NewSet(
	CoreSet,
	InvoiceFeatureSet,
	MovementFeatureSet,
	wire.Struct(new(App), "*"),
)

func InitializeApp(configFile string) (*App, func(), error) {
	panic(wire.Build(AppSet))
}
