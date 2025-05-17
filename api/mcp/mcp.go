package mcp

import (
	"context"

	"github.com/labstack/echo/v4"
	mcpSdk "github.com/mark3labs/mcp-go/mcp"
	serverSdk "github.com/mark3labs/mcp-go/server"
	"github.com/ricardogrande-masmovil/billing-mcp/api"
	invoicesPorts "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/ports"
)

// Http Controllers
type HealthController interface {
	IsHealthy(ectx echo.Context) error
}

// MCP Controllers (tool handlers)
type InvoicesController interface {
	GetInvoice(ctx context.Context, request mcpSdk.CallToolRequest) (*mcpSdk.CallToolResult, error)
	GetInvoices(ctx context.Context, request mcpSdk.CallToolRequest) (*mcpSdk.CallToolResult, error)
}

type MCPServer struct {
	HealthController
	InvoicesController
}

func NewMCPServer(healthController HealthController, invoicesController InvoicesController) *MCPServer {
	return &MCPServer{
		HealthController:   healthController,
		InvoicesController: invoicesController,
	}
}

func Setup(e *echo.Echo, s *serverSdk.MCPServer) (err error) {
	healthCtrl := api.NewHealthController()
	invoicesCtrl := invoicesPorts.NewController()

	mcp := NewMCPServer(healthCtrl, invoicesCtrl)

	sse := serverSdk.NewSSEServer(s,
		serverSdk.WithHTTPServer(e.Server),
		serverSdk.WithUseFullURLForMessageEndpoint(true),
	)

	registerHandlers(e, sse, mcp)
	registerTools(s, mcp)

	return
}

func registerHandlers(e *echo.Echo, sse *serverSdk.SSEServer, mcp *MCPServer) {
	e.GET("/health", mcp.HealthController.IsHealthy)

	e.GET("/sse", echo.WrapHandler(sse.SSEHandler()))
	e.POST("/message", echo.WrapHandler(sse.MessageHandler()))
}

func registerTools(s *serverSdk.MCPServer, mcp *MCPServer) {
	// Register tools
	s.AddTool(invoiceTool, mcp.InvoicesController.GetInvoice)
	s.AddTool(invoicesTool, mcp.InvoicesController.GetInvoices)
}
