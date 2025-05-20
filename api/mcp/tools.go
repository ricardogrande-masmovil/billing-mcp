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
)