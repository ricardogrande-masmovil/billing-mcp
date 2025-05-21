package domain_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMovementService_CreateMovement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockMovementRepository(ctrl)
	logger := zerolog.Nop()
	service := domain.NewMovementService(logger, mockRepo)

	ctx := context.Background()
	invoiceID := uuid.New()
	amount := 100.50
	movementType := model.MovementTypeCredit
	description := "Test Credit Movement"

	// Capture the argument passed to Create, as the ID is generated within NewMovement
	mockRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, m *model.Movement) error {
		assert.Equal(t, invoiceID, m.InvoiceID)
		assert.Equal(t, amount, m.Amount)
		assert.Equal(t, movementType, m.MovementType)
		assert.Equal(t, description, m.Description)
		assert.NotEqual(t, uuid.Nil, m.MovementID)
		assert.Equal(t, model.StatusPending, m.Status) // NewMovement sets status to Pending
		return nil
	}).Times(1)

	createdMovement, err := service.CreateMovement(ctx, invoiceID, amount, movementType, description)
	assert.NoError(t, err)
	assert.NotNil(t, createdMovement)
	assert.Equal(t, invoiceID, createdMovement.InvoiceID)
	assert.Equal(t, amount, createdMovement.Amount)
	assert.Equal(t, movementType, createdMovement.MovementType)
	assert.Equal(t, description, createdMovement.Description)
	assert.NotEqual(t, uuid.Nil, createdMovement.MovementID)
	assert.Equal(t, model.StatusPending, createdMovement.Status)
}

func TestMovementService_GetMovement(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockMovementRepository(ctrl)
	logger := zerolog.Nop()
	service := domain.NewMovementService(logger, mockRepo)

	ctx := context.Background()
	movementID := uuid.New()
	now := time.Now().Truncate(time.Microsecond) // Truncate for stable comparison
	expectedMovement := &model.Movement{
		MovementID:      movementID,
		InvoiceID:       uuid.New(),
		Amount:          200.00,
		MovementType:    model.MovementTypeDebit,
		Description:     "Test Debit Movement",
		TransactionDate: now,
		Status:          model.StatusInvoiced,
	}

	mockRepo.EXPECT().GetByID(ctx, movementID).Return(expectedMovement, nil).Times(1)

	retrievedMovement, err := service.GetMovement(ctx, movementID)
	assert.NoError(t, err)
	assert.Equal(t, expectedMovement, retrievedMovement)
}

func TestMovementService_UpdateMovementStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockMovementRepository(ctrl)
	logger := zerolog.Nop()
	service := domain.NewMovementService(logger, mockRepo)

	ctx := context.Background()
	movementID := uuid.New()
	originalStatus := model.StatusPending
	newStatus := model.StatusCancelled
	now := time.Now().Truncate(time.Microsecond)

	originalMovement := &model.Movement{
		MovementID:      movementID,
		InvoiceID:       uuid.New(),
		Amount:          300.00,
		MovementType:    model.MovementTypeCredit,
		Description:     "Test Update Status",
		TransactionDate: now,
		Status:          originalStatus,
	}

	mockRepo.EXPECT().GetByID(ctx, movementID).Return(originalMovement, nil).Times(1)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).DoAndReturn(func(_ context.Context, m *model.Movement) error {
		assert.Equal(t, movementID, m.MovementID)
		assert.Equal(t, newStatus, m.Status) // Check that status is updated
		return nil
	}).Times(1)

	updatedMovement, err := service.UpdateMovementStatus(ctx, movementID, newStatus)
	assert.NoError(t, err)
	assert.NotNil(t, updatedMovement)
	assert.Equal(t, movementID, updatedMovement.MovementID)
	assert.Equal(t, newStatus, updatedMovement.Status)
}

func TestMovementService_SearchMovements(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := domain.NewMockMovementRepository(ctrl)
	logger := zerolog.Nop()
	service := domain.NewMovementService(logger, mockRepo)

	ctx := context.Background()
	invoiceUUID := uuid.New()
	statusPending := model.StatusPending
	criteria := &model.SearchCriteria{ // Pass as pointer
		InvoiceID: &invoiceUUID,   // Pass as pointer
		Status:    &statusPending, // Pass as pointer
	}
	now := time.Now().Truncate(time.Microsecond)
	expectedMovements := []*model.Movement{
		{
			MovementID:      uuid.New(),
			InvoiceID:       invoiceUUID,
			Amount:          50.25,
			MovementType:    model.MovementTypeCredit, // Assuming search might return this type
			Description:     "Search Result 1",
			TransactionDate: now,
			Status:          model.StatusPending,
		},
	}

	mockRepo.EXPECT().Search(ctx, criteria).Return(expectedMovements, nil).Times(1)

	results, err := service.SearchMovements(ctx, criteria)
	assert.NoError(t, err)
	assert.Equal(t, expectedMovements, results)
}
