package domain

import (
	"context"

	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Repository interface {
	GetInvoiceByID(id model.InvoiceID) (model.Invoice, error)
	GetInvoicesByAccountId(accountId string, criteria model.Criteria) (model.Invoices, error)
	GetInvoiceLines(ctx context.Context, id model.InvoiceID) ([]model.InvoiceLine, error)
}

type Service struct {
	repo   Repository
	logger zerolog.Logger
}

func NewService(repo Repository) Service {
	return Service{
		repo:   repo,
		logger: log.With().Str("module", "invoicesService").Logger(),
	}
}

func (s Service) GetInvoiceByID(id model.InvoiceID) (model.Invoice, error) {
	s.logger.Info().Str("id", id.String()).Msg("Fetching invoice by ID")

	invoice, err := s.repo.GetInvoiceByID(id)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to fetch invoice by ID")
		return model.Invoice{}, err
	}
	return invoice, nil
}

func (s Service) GetInvoicesByCriteria(accountId string, criteria model.Criteria) (model.Invoices, error) {
	s.logger.Info().Str("account_id", accountId).Interface("criteria", criteria).Msg("Fetching invoices by criteria")

	invoices, err := s.repo.GetInvoicesByAccountId(accountId, criteria)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to fetch invoices by criteria")
		return nil, err
	}
	return invoices, nil
}

func (s Service) GetInvoiceLines(ctx context.Context, id model.InvoiceID) ([]model.InvoiceLine, error) {
	s.logger.Info().Str("id", id.String()).Msg("Fetching invoice lines by invoice ID")

	lines, err := s.repo.GetInvoiceLines(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Str("invoice_id", id.String()).Msg("Failed to fetch invoice lines")
		return nil, err
	}
	return lines, nil
}
