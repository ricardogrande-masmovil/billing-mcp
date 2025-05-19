package persistence

import (
	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence/sql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	converter        sql.InvoiceSqlConverter
	logger           zerolog.Logger
	invoiceSqlClient sql.InvoiceSqlClient
}

func NewRepository(invoiceSqlClient sql.InvoiceSqlClient, converter sql.InvoiceSqlConverter) Repository { // Inject converter
	return Repository{
		converter:        converter, // Use injected converter
		invoiceSqlClient: invoiceSqlClient,
		logger:           log.With().Str("component", "InvoicesPersistenceRepository").Logger(),
	}
}

func (r Repository) GetInvoiceByID(id domain.InvoiceID) (invoice domain.Invoice, err error) { // Changed id to domain.InvoiceID
	r.logger.Info().Str("id", id.String()).Msg("Fetching invoice by ID") // Use id.String()

	// Assuming invoiceSqlClient.GetInvoiceByID still expects a string.
	// If it can take domain.InvoiceID directly, this conversion is not needed.
	invoiceSqlModel, err := r.invoiceSqlClient.GetInvoiceByID(id.String())
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to fetch invoice by ID")
		return
	}

	invoice, err = r.converter.ConvertInvoiceToDomain(invoiceSqlModel)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to convert invoice to domain model")
		return
	}

	// TODO: fetch invoice lines and convert them to domain model

	r.logger.Info().Str("id", id.String()).Msg("Fetched invoice by ID")
	return
}

func (r Repository) GetInvoicesByAccountId(accountId string, criteria domain.Criteria) (invoices domain.Invoices, err error) { // Renamed from GetInvoicesByAccount
	r.logger.Info().Str("account_id", accountId).Interface("criteria", criteria).Msg("Fetching invoices by criteria")

	invoiceSqlModels, err := r.invoiceSqlClient.GetInvoicesByAccountId(accountId, r.converter.ConvertCriteriaToSql(criteria))
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to fetch invoices by criteria")
		return
	}

	invoices, err = r.converter.ConvertInvoicesToDomain(invoiceSqlModels)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to convert invoices to domain model")
		return
	}

	r.logger.Info().Int("count", len(invoices)).Msg("Fetched invoices by criteria")
	return
}
