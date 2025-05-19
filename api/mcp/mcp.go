package mcp

import (
	"context"

	"github.com/labstack/echo/v4"
	mcpSdk "github.com/mark3labs/mcp-go/mcp"
	serverSdk "github.com/mark3labs/mcp-go/server"
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

func Setup(e *echo.Echo, s *serverSdk.MCPServer, mcpServer *MCPServer) (err error) {
	sse := serverSdk.NewSSEServer(s,
		serverSdk.WithHTTPServer(e.Server),
		serverSdk.WithUseFullURLForMessageEndpoint(true),
	)

	registerHandlers(e, sse, mcpServer)
	registerTools(s, mcpServer)

	return
}

func registerHandlers(e *echo.Echo, sse *serverSdk.SSEServer, mcp *MCPServer) {
	e.GET("/health", mcp.HealthController.IsHealthy)

	e.GET("/sse", echo.WrapHandler(sse.SSEHandler()))
	e.POST("/message", echo.WrapHandler(sse.MessageHandler()))
}

func registerTools(s *serverSdk.MCPServer, mcp *MCPServer) {
	s.AddTool(invoiceTool, mcp.InvoicesController.GetInvoice)
	s.AddTool(invoicesTool, mcp.InvoicesController.GetInvoices)
}
