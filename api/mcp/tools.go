package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
)

var (
	invoiceTool = mcp.NewTool(
		"GetInvoice",
		mcp.WithDescription("Get an invoice by ID"),
		mcp.WithString("accountId", mcp.Required(), mcp.Description("The ID of the account to retrieve the invoice for")),
		mcp.WithString("invoiceId", mcp.Required(), mcp.Description("The ID of the invoice to retrieve")),
	)

	invoicesTool = mcp.NewTool(
		"GetInvoices",
		mcp.WithDescription("Get all invoices for an account"),
		mcp.WithString("accountId", mcp.Required(), mcp.Description("The ID of the account to retrieve invoices for")),
		mcp.WithString("status", mcp.Description("The status of the invoices to retrieve")),
		mcp.WithString("issueDateFrom", mcp.Description("The start date of the invoices to retrieve in RFC3339 format")),
		mcp.WithString("issueDateTo", mcp.Description("The end date of the invoices to retrieve in RFC3339 format")),
	)

	invoiceMovementsTool = mcp.NewTool(
		"GetInvoiceMovements",
		mcp.WithDescription("Get all movements/lines for a specific invoice"),
		mcp.WithString("accountId", mcp.Required(), mcp.Description("The ID of the account")),
		mcp.WithString("invoiceId", mcp.Required(), mcp.Description("The ID of the invoice to retrieve movements for")),
	)

	movementTool = mcp.NewTool(
		"GetMovement",
		mcp.WithDescription("Get a specific movement by ID"),
		mcp.WithString("accountId", mcp.Required(), mcp.Description("The ID of the account")),
		mcp.WithString("movementId", mcp.Required(), mcp.Description("The ID of the movement to retrieve")),
	)
)
