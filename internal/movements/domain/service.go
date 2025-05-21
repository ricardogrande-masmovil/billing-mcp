package domain

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/rs/zerolog"
)

// MovementRepository defines the interface for movement persistence.
// It's implemented by an adapter in the infrastructure layer.
type MovementRepository interface {
	Create(ctx context.Context, movement *model.Movement) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Movement, error)
	Update(ctx context.Context, movement *model.Movement) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, criteria *model.SearchCriteria) ([]*model.Movement, error)
}

// MovementService provides business logic for movements.
// It depends on the MovementRepository for data access.
type MovementService struct {
	logger     zerolog.Logger
	repository MovementRepository
}

// NewMovementService creates a new MovementService.
func NewMovementService(logger zerolog.Logger, repository MovementRepository) *MovementService {
	return &MovementService{
		logger:     logger.With().Str("service", "MovementService").Logger(),
		repository: repository,
	}
}

// CreateMovement creates a new movement.
func (s *MovementService) CreateMovement(ctx context.Context, invoiceID uuid.UUID, amount float64, movementType model.MovementType, description string) (*model.Movement, error) {
	log := s.logger.With().Str("method", "CreateMovement").Logger()

	movement, err := model.NewMovement(invoiceID, amount, movementType, description)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create new movement domain model")
		return nil, fmt.Errorf("failed to create new movement: %w", err)
	}

	if err := s.repository.Create(ctx, movement); err != nil {
		log.Error().Err(err).Msg("Failed to save movement to repository")
		return nil, fmt.Errorf("failed to save movement: %w", err)
	}

	log.Info().Str("movementID", movement.MovementID.String()).Msg("Movement created successfully")
	return movement, nil
}

// GetMovement retrieves a movement by its ID.
func (s *MovementService) GetMovement(ctx context.Context, id uuid.UUID) (*model.Movement, error) {
	log := s.logger.With().Str("method", "GetMovement").Str("movementID", id.String()).Logger()

	movement, err := s.repository.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movement from repository")
		return nil, fmt.Errorf("failed to get movement with ID %s: %w", id, err)
	}

	log.Info().Msg("Movement retrieved successfully")
	return movement, nil
}

// UpdateMovementStatus updates the status of an existing movement.
func (s *MovementService) UpdateMovementStatus(ctx context.Context, id uuid.UUID, status model.Status) (*model.Movement, error) {
	log := s.logger.With().Str("method", "UpdateMovementStatus").Str("movementID", id.String()).Str("newStatus", string(status)).Logger()

	movement, err := s.repository.GetByID(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get movement for status update")
		return nil, fmt.Errorf("failed to get movement with ID %s for update: %w", id, err)
	}

	// TODO: Add business logic for status transitions if needed
	movement.Status = status

	if err := s.repository.Update(ctx, movement); err != nil {
		log.Error().Err(err).Msg("Failed to update movement status in repository")
		return nil, fmt.Errorf("failed to update movement status: %w", err)
	}

	log.Info().Msg("Movement status updated successfully")
	return movement, nil
}

// SearchMovements searches for movements based on criteria.
func (s *MovementService) SearchMovements(ctx context.Context, criteria *model.SearchCriteria) ([]*model.Movement, error) {
	log := s.logger.With().Str("method", "SearchMovements").Logger()

	movements, err := s.repository.Search(ctx, criteria)
	if err != nil {
		log.Error().Err(err).Msg("Failed to search movements in repository")
		return nil, fmt.Errorf("failed to search movements: %w", err)
	}

	log.Info().Int("count", len(movements)).Msg("Movements searched successfully")
	return movements, nil
}
