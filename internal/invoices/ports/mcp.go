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
	GetInvoiceLines(ctx context.Context, id domain.InvoiceID) ([]domain.InvoiceLine, error)
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

	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		c.logger.Error().Msg("Arguments must be a map[string]interface{}")
		return mcp.NewToolResultErrorFromErr("Invalid arguments type", errors.New("arguments must be a map[string]interface{}")), nil
	}
	requestedInvoiceId, ok := args["invoiceId"].(string)
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

	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		c.logger.Error().Msg("Arguments must be a map[string]interface{}")
		return mcp.NewToolResultErrorFromErr("Invalid arguments type", errors.New("arguments must be a map[string]interface{}")), nil
	}
	accountId, ok := args["accountId"].(string)
	if !ok || accountId == "" {
		c.logger.Error().Msg("Account ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingAccountId), nil
	}

	criteria, err := c.converter.ConvertRequestArgsToCriteria(args)

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

	// Extract arguments as map[string]interface{}
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		c.logger.Error().Msg("Arguments must be a map[string]interface{}")
		return mcp.NewToolResultErrorFromErr("Invalid arguments type", errors.New("arguments must be a map[string]interface{}")), nil
	}

	// Extract invoiceId from arguments
	requestedInvoiceId, ok := args["invoiceId"].(string)
	if !ok || requestedInvoiceId == "" {
		c.logger.Error().Msg("Invoice ID is required")
		return mcp.NewToolResultErrorFromErr("Missing request parameter", ErrMissingInvoiceId), nil
	}

	// Extract accountId from arguments (not used in this implementation but required by API contract)
	_, ok = args["accountId"].(string)
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

	// Fetch invoice lines from service
	lines, err := c.service.GetInvoiceLines(ctx, invoiceId)
	if err != nil {
		c.logger.Error().Err(err).Str("invoiceId", requestedInvoiceId).Msg("Failed to get invoice lines")
		return mcp.NewToolResultErrorFromErr("Failed to retrieve invoice lines", err), nil
	}

	// Convert lines to InvoiceMovementDTO
	var movementDTOs []InvoiceMovementDTO
	for _, line := range lines {
		movementDTOs = append(movementDTOs, InvoiceMovementDTO{
			MovementID:       line.MovementID.String(),
			Description:      line.Description,
			Amount:           line.AmountWithTax,
			AmountWithoutTax: line.AmountWithoutTax,
			AmountWithTax:    line.AmountWithTax,
			TaxPercentage:    line.TaxPercentage,
			OperationType:    line.OperationType,
		})
	}

	// Convert movements to JSON
	jsonData, err := c.converter.ConvertInvoiceMovementsToJson(movementDTOs)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to convert invoice movements to JSON")
		return nil, err
	}

	c.logger.Info().Str("invoiceId", requestedInvoiceId).Int("movementsCount", len(movementDTOs)).Msg("Successfully retrieved invoice movements")
	response := mcp.NewToolResultText(string(jsonData))
	return response, nil
}
