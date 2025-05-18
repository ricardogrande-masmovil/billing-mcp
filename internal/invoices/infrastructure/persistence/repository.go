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

func NewRepository(invoiceSqlClient sql.InvoiceSqlClient) Repository {
	return Repository{
		converter:        sql.NewInvoiceSqlConverter(),
		invoiceSqlClient: invoiceSqlClient,
		logger:           log.With().Str("component", "InvoicesPersistenceRepository").Logger(),
	}
}

func (r Repository) GetInvoiceByID(id string) (invoice domain.Invoice, err error) {
	r.logger.Info().Str("id", id).Msg("Fetching invoice by ID")

	invoiceSqlModel, err := r.invoiceSqlClient.GetInvoiceByID(id)
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

	r.logger.Info().Str("id", id).Msg("Fetched invoice by ID")
	return
}

func (r Repository) GetInvoicesByAccount(accountId string, criteria domain.Criteria) (invoices domain.Invoices, err error) {
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
