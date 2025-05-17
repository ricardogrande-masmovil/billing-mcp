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
	)
)