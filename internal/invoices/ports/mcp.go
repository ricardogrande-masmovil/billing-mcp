package ports

import (
	"context"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type controller struct {
	converter Converter
	logger    zerolog.Logger
}

func NewController() controller {
	return controller{
		converter: NewConverter(),
		logger:    log.With().Str("module", "invoicesMcpController").Logger(),
	}
}

func (c controller) GetInvoice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c.logger.Info().Msg("Processing request in GetInvoice tool")

	mockInvoice := domain.Invoice{
		ID:                    domain.InvoiceID("invoice-123"),
		TotalAmountWithoutTax: 100,
		TotalAmountWithTax:    120,
		Status:                domain.InvoiceStatusPaid,
		IssueDate:             time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		DueDate:               time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
	}

	jsonData, err := c.converter.ConvertDomainInvoiceToJsonInvoice(mockInvoice)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoice to JSON")
		return nil, err
	}

	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}

func (c controller) GetInvoices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c.logger.Info().Msg("Processing request in GetInvoices tool")

	mockInvoices := domain.Invoices{
		domain.Invoice{
			ID:                    domain.InvoiceID("invoice-123"),
			TotalAmountWithoutTax: 100,
			TotalAmountWithTax:    120,
			Status:                domain.InvoiceStatusPaid,
			IssueDate:             time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			DueDate:               time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
		},
		domain.Invoice{
			ID:                    domain.InvoiceID("invoice-456"),
			TotalAmountWithoutTax: 200,
			TotalAmountWithTax:    240,
			Status:                domain.InvoiceStatusUnpaid,
			IssueDate:             time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			DueDate:               time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	jsonData, err := c.converter.ConvertDomainInvoicesToJsonInvoices(mockInvoices)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoices to JSON")
		return nil, err
	}

	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}
