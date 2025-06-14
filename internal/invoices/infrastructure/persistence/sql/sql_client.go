package sql

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type InvoiceSqlClient struct {
	db         *gorm.DB
	maxRetries int
	logger     zerolog.Logger
}

func NewInvoiceSqlClient(db *gorm.DB, maxRetries int) InvoiceSqlClient {
	return InvoiceSqlClient{
		db:         db,
		maxRetries: maxRetries,
		logger:     log.With().Str("component", "InvoicesPersistenceRepository").Logger(),
	}
}

func (c InvoiceSqlClient) GetInvoiceByID(id string) (invoice Invoice, err error) {
	c.logger.Info().Str("id", id).Msg("Fetching invoice by ID")

	queryFn := func() *gorm.DB {
		return c.db.Where("id = ?", id).First(&invoice)
	}

	rowsAffected, err := c.RunWithRetry(queryFn, c.maxRetries)
	if err != nil {
		return
	}

	c.logger.Info().Int("rows_affected", rowsAffected).Msg("Fetched invoice by ID")
	return
}

func (c InvoiceSqlClient) GetInvoicesByAccountId(accountId string, criteria map[string]interface{}) (invoices []Invoice, err error) {
	c.logger.Info().Interface("criteria", criteria).Msg("Fetching invoices by criteria")

	queryFn := func() *gorm.DB {
		query := c.db.Where("account_id = ?", accountId)
		if criteria["status"] != nil {
			query = query.Where("status = ?", criteria["status"])
		}
		if criteria["issue_date_from"] != nil {
			query = query.Where("issue_date >= ?", criteria["issue_date_from"])
		}
		if criteria["issue_date_to"] != nil {
			query = query.Where("issue_date <= ?", criteria["issue_date_to"])
		}
		query = query.Order("issue_date DESC")
		return query.Find(&invoices)
	}

	rowsAffected, err := c.RunWithRetry(queryFn, c.maxRetries)
	if err != nil {
		return
	}

	c.logger.Info().Int("rows_affected", rowsAffected).Msg("Fetched invoices by criteria")
	return
}

// GetInvoiceLinesByInvoiceID retrieves all invoice lines for a specific invoice
func (c InvoiceSqlClient) GetInvoiceLinesByInvoiceID(ctx context.Context, invoiceID string) ([]InvoiceLine, error) {
	c.logger.Info().Str("invoice_id", invoiceID).Msg("Fetching invoice lines by invoice ID")

	var lines []InvoiceLine

	queryFn := func() *gorm.DB {
		return c.db.WithContext(ctx).
			Where("invoice_id = ?", invoiceID).
			Find(&lines)
	}

	_, err := c.RunWithRetry(queryFn, c.maxRetries)
	if err != nil {
		c.logger.Error().Err(err).Str("invoice_id", invoiceID).Msg("Failed to fetch invoice lines")
		return nil, err
	}

	c.logger.Info().Str("invoice_id", invoiceID).Int("count", len(lines)).Msg("Successfully fetched invoice lines")
	return lines, nil
}

func (c InvoiceSqlClient) RunWithRetry(queryFn func() *gorm.DB, retries int) (rowsAffected int, err error) {
	for i := 0; i < retries; i++ {
		result := queryFn()
		err = result.Error
		if err == nil {
			return int(result.RowsAffected), nil
		}
		c.logger.Error().Err(err).Msg("Query failed, retrying...")
	}

	return 0, err
}
