package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	sqlmodel "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/infrastructure/persistence/sql"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// MovementSQLRepository is a GORM-based implementation of the MovementRepository interface.
// It handles database operations for movements.
type MovementSQLRepository struct {
	db     *gorm.DB
	logger zerolog.Logger
}

// NewMovementSQLRepository creates a new MovementSQLRepository.
func NewMovementSQLRepository(db *gorm.DB, logger zerolog.Logger) domain.MovementRepository {
	return &MovementSQLRepository{
		db:     db,
		logger: logger.With().Str("repository", "MovementSQLRepository").Logger(),
	}
}

// Create saves a new movement to the database.
func (r *MovementSQLRepository) Create(ctx context.Context, movement *model.Movement) error {
	log := r.logger.With().Str("method", "Create").Str("movementID", movement.MovementID.String()).Logger()
	sqlMovement := sqlmodel.ToSQLMovement(movement)

	if err := r.db.WithContext(ctx).Create(sqlMovement).Error; err != nil {
		log.Error().Err(err).Msg("Failed to create movement in database")
		return fmt.Errorf("failed to create movement: %w", err)
	}
	log.Info().Msg("Movement created successfully in database")
	return nil
}

// GetByID retrieves a movement by its ID from the database.
func (r *MovementSQLRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Movement, error) {
	log := r.logger.With().Str("method", "GetByID").Str("movementID", id.String()).Logger()
	var sqlMovement sqlmodel.Movement

	if err := r.db.WithContext(ctx).First(&sqlMovement, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn().Msg("Movement not found in database")
			return nil, fmt.Errorf("movement with ID %s not found: %w", id, err) // Consider a domain-specific error
		}
		log.Error().Err(err).Msg("Failed to get movement from database")
		return nil, fmt.Errorf("failed to get movement with ID %s: %w", id, err)
	}

	log.Info().Msg("Movement retrieved successfully from database")
	return sqlmodel.ToDomainMovement(&sqlMovement), nil
}

// Update modifies an existing movement in the database.
func (r *MovementSQLRepository) Update(ctx context.Context, movement *model.Movement) error {
	log := r.logger.With().Str("method", "Update").Str("movementID", movement.MovementID.String()).Logger()
	sqlMovement := sqlmodel.ToSQLMovement(movement)

	// Ensure the ID is set for the update operation
	if sqlMovement.ID == uuid.Nil {
		log.Error().Msg("Cannot update movement with nil ID")
		return fmt.Errorf("cannot update movement: ID is nil")
	}

	result := r.db.WithContext(ctx).Model(&sqlmodel.Movement{}).Where("id = ?", sqlMovement.ID).Updates(sqlMovement)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to update movement in database")
		return fmt.Errorf("failed to update movement: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		log.Warn().Msg("Movement not found for update or no changes made")
		// Consider returning a specific error if the record was not found vs. no changes made
		return fmt.Errorf("movement with ID %s not found for update or no changes made", movement.MovementID)
	}

	log.Info().Msg("Movement updated successfully in database")
	return nil
}

// Delete removes a movement from the database (soft delete if BaseModel is used).
func (r *MovementSQLRepository) Delete(ctx context.Context, id uuid.UUID) error {
	log := r.logger.With().Str("method", "Delete").Str("movementID", id.String()).Logger()

	result := r.db.WithContext(ctx).Delete(&sqlmodel.Movement{}, "id = ?", id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Failed to delete movement from database")
		return fmt.Errorf("failed to delete movement with ID %s: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		log.Warn().Msg("Movement not found for deletion")
		return fmt.Errorf("movement with ID %s not found for deletion", id)
	}

	log.Info().Msg("Movement deleted successfully from database")
	return nil
}

// Search retrieves movements based on the provided criteria.
func (r *MovementSQLRepository) Search(ctx context.Context, criteria *model.SearchCriteria) ([]*model.Movement, error) {
	log := r.logger.With().Str("method", "Search").Logger()
	var sqlMovements []*sqlmodel.Movement

	dbQuery := r.db.WithContext(ctx).Model(&sqlmodel.Movement{})

	if criteria.InvoiceID != nil {
		dbQuery = dbQuery.Where("invoice_id = ?", *criteria.InvoiceID)
	}
	if criteria.Status != nil {
		dbQuery = dbQuery.Where("status = ?", *criteria.Status)
	}
	// Add other criteria filters here

	if err := dbQuery.Find(&sqlMovements).Error; err != nil {
		log.Error().Err(err).Msg("Failed to search movements in database")
		return nil, fmt.Errorf("failed to search movements: %w", err)
	}

	log.Info().Int("count", len(sqlMovements)).Msg("Movements searched successfully in database")
	return sqlmodel.ToDomainMovements(sqlMovements), nil
}
