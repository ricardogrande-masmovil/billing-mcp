package ports

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	mcpSdk "github.com/mark3labs/mcp-go/mcp"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/rs/zerolog"
)

// MCPMovementsHandler handles MCP requests for movements
type MCPMovementsHandler struct {
	movementService domain.MovementService
	logger          zerolog.Logger
}

// NewMCPMovementsHandler creates a new MCPMovementsHandler
func NewMCPMovementsHandler(movementService domain.MovementService, logger zerolog.Logger) *MCPMovementsHandler {
	return &MCPMovementsHandler{
		movementService: movementService,
		logger:          logger.With().Str("component", "MCPMovementsHandler").Logger(),
	}
}

// GetMovement handles the GetMovement MCP tool
func (h *MCPMovementsHandler) GetMovement(ctx context.Context, request mcpSdk.CallToolRequest) (*mcpSdk.CallToolResult, error) {
	log := h.logger.With().Str("method", "GetMovement").Logger()
	log.Debug().Msg("Processing GetMovement request")

	// Extract arguments map from request
	args, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		log.Error().Msg("Invalid arguments type in request")
		return mcpSdk.NewToolResultErrorFromErr("Invalid arguments", fmt.Errorf("arguments must be a map")), nil
	}

	// Extract accountId from arguments (we validate it but don't use it in this implementation)
	accountIDStr, ok := args["accountId"].(string)
	if !ok || accountIDStr == "" {
		log.Error().Msg("Missing or invalid accountId parameter")
		return mcpSdk.NewToolResultErrorFromErr("Missing parameter", fmt.Errorf("accountId is required")), nil
	}

	// Extract movementId from arguments
	movementIDStr, ok := args["movementId"].(string)
	if !ok || movementIDStr == "" {
		log.Error().Msg("Missing or invalid movementId parameter")
		return mcpSdk.NewToolResultErrorFromErr("Missing parameter", fmt.Errorf("movementId is required")), nil
	}

	// Parse the UUID
	movementID, err := uuid.Parse(movementIDStr)
	if err != nil {
		log.Error().Err(err).Str("movementId", movementIDStr).Msg("Failed to parse movementId")
		return mcpSdk.NewToolResultErrorFromErr("Invalid format", fmt.Errorf("invalid movement ID format: %w", err)), nil
	}

	// Fetch the movement from domain service
	movement, err := h.movementService.GetMovement(ctx, movementID)
	if err != nil {
		log.Error().Err(err).Str("movementId", movementIDStr).Msg("Failed to get movement")
		return nil, fmt.Errorf("failed to retrieve movement: %w", err)
	}

	// Convert to response DTO
	response := convertToMovementDTO(movement)

	// Convert to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		log.Error().Err(err).Msg("Failed to convert movement to JSON")
		return nil, fmt.Errorf("failed to convert movement to JSON: %w", err)
	}

	log.Info().Str("movementId", movementIDStr).Msg("Successfully retrieved movement")

	// Return the response as text
	result := mcpSdk.NewToolResultText(string(jsonData))
	return result, nil
}

// Helper functions for conversion

// convertToMovementDTO converts a domain Movement to a DTO
func convertToMovementDTO(m *model.Movement) *MovementDTO {
	return &MovementDTO{
		ID:              m.MovementID.String(),
		InvoiceID:       m.InvoiceID.String(),
		Amount:          m.Amount,
		MovementType:    string(m.MovementType),
		Description:     m.Description,
		TransactionDate: m.TransactionDate.Format(time.RFC3339),
		Status:          string(m.Status),
	}
}
