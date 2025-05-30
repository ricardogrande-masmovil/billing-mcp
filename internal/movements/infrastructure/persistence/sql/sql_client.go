package sql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MovementSqlClient handles database operations for movements.
type MovementSqlClient struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewMovementSqlClient creates a new MovementSqlClient.
func NewMovementSqlClient(db *gorm.DB, logger zerolog.Logger) *MovementSqlClient {
	return &MovementSqlClient{
		db:     db,
		logger: logger.With().Str("component", "MovementSqlClient").Logger(),
	}
}

// CreateMovement creates a new movement record in the database.
func (c *MovementSqlClient) CreateMovement(ctx context.Context, m *Movement) error {
	log := c.logger.With().Str("method", "CreateMovement").Logger()
	log.Debug().Interface("movement", m).Msg("Creating movement")

	if err := c.db.WithContext(ctx).Create(m).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create movement")
		return fmt.Errorf("failed to create movement: %w", err)
	}
	log.Info().Msg("Movement created successfully")
	return nil
}

// GetMovementByID retrieves a movement by its ID.
func (c *MovementSqlClient) GetMovementByID(ctx context.Context, id uuid.UUID) (*Movement, error) {
	log := c.logger.With().Str("method", "GetMovementByID").Stringer("movementID", id).Logger()
	log.Debug().Msg("Getting movement by ID")

	var movement Movement
	if err := c.db.WithContext(ctx).First(&movement, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn().Msg("Movement not found")
			return nil, fmt.Errorf("movement with ID %s not found: %w", id, gorm.ErrRecordNotFound)
		}
		log.Error().Err(err).Msg("Failed to get movement by ID")
		return nil, fmt.Errorf("failed to get movement by ID %s: %w", id, err)
	}
	log.Info().Msg("Movement retrieved successfully")
	return &movement, nil
}

// UpdateMovement updates an existing movement in the database.
// It specifically updates the Status.
func (c *MovementSqlClient) UpdateMovement(ctx context.Context, m *Movement) error {
	log := c.logger.With().Str("method", "UpdateMovement").Stringer("movementID", m.ID).Logger()
	log.Debug().Interface("movement", m).Msg("Updating movement")

	result := c.db.WithContext(ctx).Model(&Movement{}).Where("id = ?", m.ID).Updates(m)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to update movement")
		return fmt.Errorf("failed to update movement with ID %s: %w", m.ID, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Warn().Msg("Movement not found for update or no changes made")
		return fmt.Errorf("movement with ID %s not found for update or no changes made: %w", m.ID, gorm.ErrRecordNotFound)
	}
	log.Info().Msg("Movement updated successfully")
	return nil
}

// DeleteMovement deletes a movement by its ID.
func (c *MovementSqlClient) DeleteMovement(ctx context.Context, id uuid.UUID) error {
	log := c.logger.With().Str("method", "DeleteMovement").Stringer("movementID", id).Logger()
	log.Debug().Msg("Deleting movement")

	result := c.db.WithContext(ctx).Delete(&Movement{}, "id = ?", id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to delete movement")
		return fmt.Errorf("failed to delete movement with ID %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Warn().Msg("Movement not found for deletion")
		return fmt.Errorf("movement with ID %s not found for deletion: %w", id, gorm.ErrRecordNotFound)
	}
	log.Info().Msg("Movement deleted successfully")
	return nil
}

// SearchMovements searches for movements based on criteria.
func (c *MovementSqlClient) SearchMovements(ctx context.Context, criteria *model.SearchCriteria) ([]Movement, error) {
	log := c.logger.With().Str("method", "SearchMovements").Interface("criteria", criteria).Logger()
	log.Debug().Msg("Searching movements")

	var movements []Movement
	query := c.db.WithContext(ctx)

	if criteria.InvoiceID != nil {
		query = query.Where("invoice_id = ?", *criteria.InvoiceID)
	}
	if criteria.Status != nil {
		query = query.Where("status = ?", criteria.Status.String())
	}
	// Add other criteria as needed, e.g., date ranges, movement type

	if err := query.Order("transaction_date DESC").Find(&movements).Error; err != nil {
		log.Error().Err(err).Msg("Failed to search movements")
		return nil, fmt.Errorf("failed to search movements: %w", err)
	}

	log.Info().Int("count", len(movements)).Msg("Movements search completed")
	return movements, nil
}
