package persistence

import (
	"errors"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ZerologGormWriter is a custom writer to bridge GORM's logger with zerolog.
type ZerologGormWriter struct {
	logger zerolog.Logger
}

// NewZerologGormWriter creates a new writer for GORM.
func NewZerologGormWriter(logger zerolog.Logger) *ZerologGormWriter {
	return &ZerologGormWriter{logger: logger}
}

// Printf implements gormlogger.Writer interface
func (z *ZerologGormWriter) Printf(format string, args ...interface{}) {
	z.logger.Printf(format, args...)
}

// NewSqlClient initializes and returns a new Gorm DB instance for PostgreSQL.
// dsn is the Data Source Name for connecting to the PostgreSQL database.
// Example DSN: "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable TimeZone=Asia/Shanghai"
func NewSqlClient(dsn string) (*gorm.DB, error) {
	gormLog := log.With().Str("component", "gorm").Logger()

	gormDBLogger := gormlogger.New(
		NewZerologGormWriter(gormLog),
		gormlogger.Config{
			SlowThreshold:             200,
			LogLevel:                  gormlogger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormDBLogger,
	})

	if err != nil {
		gormLog.Error().Err(err).Msg("Failed to connect to PostgreSQL database")
		return nil, errors.New("failed to connect to PostgreSQL database")
	}

	gormLog.Info().Msg("PostgreSQL database connection established.")
	return db, nil
}
