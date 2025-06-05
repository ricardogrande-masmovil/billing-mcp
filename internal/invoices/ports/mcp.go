package ports

import (
	"context"
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	ErrMissingInvoiceId = errors.New("invoice_id is required")
	ErrMissingAccountId = errors.New("account_id is required")
)

type InvoiceService interface {
	GetInvoiceByID(id domain.InvoiceID) (domain.Invoice, error)
	GetInvoicesByCriteria(accountId string, criteria domain.Criteria) (domain.Invoices, error)
	GetInvoiceMovements(id domain.InvoiceID) ([]InvoiceMovementDTO, error)
}

type controller struct {
	service   InvoiceService
	converter Converter
	logger    zerolog.Logger
}

func NewController(service InvoiceService) controller {
	return controller{
		service:   service,
		converter: NewConverter(),
		logger:    log.With().Str("module", "invoicesMcpController").Logger(),
	}
}

func (c controller) GetInvoice(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c.logger.Info().Msg("Processing request in GetInvoice tool")

	requestedInvoiceId, ok := request.Params.Arguments["invoiceId"].(string)
	if !ok || requestedInvoiceId == "" {
		c.logger.Error().Msg("Invoice ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingInvoiceId), nil
	}

	invoiceId, err := domain.ParseInvoiceID(requestedInvoiceId)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to parse invoice ID")
		return mcp.NewToolResultErrorFromErr("Invalid invoice ID format", err), nil
	}

	invoice, err := c.service.GetInvoiceByID(invoiceId)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to fetch invoice by ID")
		return nil, err
	}

	jsonData, err := c.converter.ConvertDomainInvoiceToJsonInvoice(invoice)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoice to JSON")
		return nil, err
	}

	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}

func (c controller) GetInvoices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c.logger.Info().Msg("Processing request in GetInvoices tool")

	accountId, ok := request.Params.Arguments["accountId"].(string)
	if !ok || accountId == "" {
		c.logger.Error().Msg("Account ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingAccountId), nil
	}

	criteria, err := c.converter.ConvertRequestArgsToCriteria(request.Params.Arguments)

	invoices, err := c.service.GetInvoicesByCriteria(accountId, criteria)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to fetch invoices by criteria")
		return nil, err
	}

	jsonData, err := c.converter.ConvertDomainInvoicesToJsonInvoices(invoices)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoices to JSON")
		return nil, err
	}

	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}

func (c controller) GetInvoiceMovements(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	c.logger.Info().Msg("Processing request in GetInvoiceMovements tool")

	// Extract invoiceId from request
	requestedInvoiceId, ok := request.Params.Arguments["invoiceId"].(string)
	if !ok || requestedInvoiceId == "" {
		c.logger.Error().Msg("Invoice ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingInvoiceId), nil
	}

	// Extract accountId from request (not used in this implementation but required by API contract)
	_, ok = request.Params.Arguments["accountId"].(string)
	if !ok {
		c.logger.Error().Msg("Account ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingAccountId), nil
	}

	// Parse the invoice ID
	invoiceId, err := domain.ParseInvoiceID(requestedInvoiceId)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to parse invoice ID")
		return mcp.NewToolResultErrorFromErr("Invalid invoice ID format", err), nil
	}

	// Fetch invoice movements from service
	movements, err := c.service.GetInvoiceMovements(invoiceId)
	if err != nil {
		c.logger.Error().Err(err).Str("invoiceId", requestedInvoiceId).Msg("Failed to get invoice movements")
		return mcp.NewToolResultErrorFromErr("Failed to retrieve invoice movements", err), nil
	}

	// Convert movements to JSON
	jsonData, err := c.converter.ConvertInvoiceMovementsToJson(movements)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoice movements to JSON")
		return nil, err
	}

	c.logger.Info().Str("invoiceId", requestedInvoiceId).Int("movementsCount", len(movements)).Msg("Successfully retrieved invoice movements")
	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}
