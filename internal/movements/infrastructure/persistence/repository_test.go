package persistence_test

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/infrastructure/persistence"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestMovementSQLRepository_Create(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New() // This ID will be part of the SQL model
	invoiceID := uuid.New()
	transactionDate := time.Now().UTC().Truncate(time.Microsecond)
	mov := &model.Movement{
		MovementID:      movementID,
		InvoiceID:       invoiceID,
		Amount:          150.75,
		MovementType:    model.MovementTypeDebit,
		Description:     "Test Debit for Repo",
		TransactionDate: transactionDate,
		Status:          model.StatusPending,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "movements" ("id","invoice_id","amount","movement_type","description","transaction_date","status","created_at","updated_at") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING "id"`)).
		WithArgs(movementID, mov.InvoiceID, mov.Amount, string(mov.MovementType), mov.Description, mov.TransactionDate, string(mov.Status), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(movementID)) // GORM expects the ID back
	mock.ExpectCommit()

	err := repo.Create(ctx, mov)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_GetByID_Found(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New()
	invoiceID := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)

	rows := sqlmock.NewRows([]string{"id", "invoice_id", "amount", "movement_type", "description", "transaction_date", "status", "created_at", "updated_at"}).
		AddRow(movementID, invoiceID, 250.00, "CREDIT", "Repo Get Test", now, "INVOICED", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements" WHERE id = $1 ORDER BY "movements"."id" LIMIT 1`)).
		WithArgs(movementID).
		WillReturnRows(rows)

	retrievedMovement, err := repo.GetByID(ctx, movementID)

	assert.NoError(t, err)
	require.NotNil(t, retrievedMovement)
	assert.Equal(t, movementID, retrievedMovement.MovementID)
	assert.Equal(t, invoiceID, retrievedMovement.InvoiceID)
	assert.Equal(t, 250.00, retrievedMovement.Amount)
	assert.Equal(t, model.MovementTypeCredit, retrievedMovement.MovementType)
	assert.Equal(t, "Repo Get Test", retrievedMovement.Description)
	assert.True(t, now.Equal(retrievedMovement.TransactionDate)) // Use Equal for time comparison
	assert.Equal(t, model.StatusInvoiced, retrievedMovement.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_GetByID_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements" WHERE id = $1 ORDER BY "movements"."id" LIMIT 1`)).
		WithArgs(movementID).
		WillReturnError(gorm.ErrRecordNotFound)

	retrievedMovement, err := repo.GetByID(ctx, movementID)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "not found")) // Check for gorm.ErrRecordNotFound or a wrapped version
	assert.Nil(t, retrievedMovement)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Update_Status(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New()
	invoiceID := uuid.New()
	now := time.Now().UTC().Truncate(time.Microsecond)
	newStatus := model.StatusInvoiced

	movToUpdate := &model.Movement{
		MovementID:      movementID,
		InvoiceID:       invoiceID,
		Amount:          300.50,
		MovementType:    model.MovementTypeCredit,
		Description:     "Original Description",
		TransactionDate: now,       // This usually doesn't change on status update
		Status:          newStatus, // Update the status in the model to be passed to Update
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE "movements" SET "invoice_id"=$1,"amount"=$2,"movement_type"=$3,"description"=$4,"transaction_date"=$5,"status"=$6,"updated_at"=$7 WHERE "id" = $8`)).
		WithArgs(movToUpdate.InvoiceID, movToUpdate.Amount, string(movToUpdate.MovementType), movToUpdate.Description, movToUpdate.TransactionDate, string(movToUpdate.Status), sqlmock.AnyArg(), movementID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	err := repo.Update(ctx, movToUpdate)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Delete_Found(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "movements" WHERE id = $1`)).
		WithArgs(movementID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, movementID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Delete_NotFound(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	movementID := uuid.New()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM "movements" WHERE id = $1`)).
		WithArgs(movementID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected
	mock.ExpectCommit()

	err := repo.Delete(ctx, movementID)

	assert.Error(t, err) // Expect an error because no rows were affected
	assert.True(t, strings.Contains(err.Error(), "not found for deletion"))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Search_ByInvoiceID(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	searchInvoiceID := uuid.New()
	criteria := &model.SearchCriteria{
		InvoiceID: &searchInvoiceID,
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	rows := sqlmock.NewRows([]string{"id", "invoice_id", "amount", "movement_type", "description", "transaction_date", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), searchInvoiceID, 10.0, "DEBIT", "Search 1", now, "PENDING", now, now).
		AddRow(uuid.New(), searchInvoiceID, 20.0, "CREDIT", "Search 2", now, "INVOICED", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements" WHERE invoice_id = $1`)).
		WithArgs(searchInvoiceID).
		WillReturnRows(rows)

	results, err := repo.Search(ctx, criteria)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	if len(results) == 2 {
		assert.Equal(t, searchInvoiceID, results[0].InvoiceID)
		assert.Equal(t, 10.0, results[0].Amount)
		assert.Equal(t, searchInvoiceID, results[1].InvoiceID)
		assert.Equal(t, 20.0, results[1].Amount)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Search_ByStatus(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	searchStatus := model.StatusInvoiced
	criteria := &model.SearchCriteria{
		Status: &searchStatus,
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	movementID1 := uuid.New()
	invoiceID1 := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "invoice_id", "amount", "movement_type", "description", "transaction_date", "status", "created_at", "updated_at"}).
		AddRow(movementID1, invoiceID1, 30.0, "CREDIT", "Search Status Type", now, "INVOICED", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements" WHERE status = $1`)).
		WithArgs(string(searchStatus)). // GORM will convert enum to its underlying type (string)
		WillReturnRows(rows)

	results, err := repo.Search(ctx, criteria)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	if len(results) == 1 {
		assert.Equal(t, movementID1, results[0].MovementID)
		assert.Equal(t, model.StatusInvoiced, results[0].Status)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Search_ByInvoiceIDAndStatus(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	searchInvoiceID := uuid.New()
	searchStatus := model.StatusPending
	criteria := &model.SearchCriteria{
		InvoiceID: &searchInvoiceID,
		Status:    &searchStatus,
	}

	now := time.Now().UTC().Truncate(time.Microsecond)
	movementID1 := uuid.New()
	rows := sqlmock.NewRows([]string{"id", "invoice_id", "amount", "movement_type", "description", "transaction_date", "status", "created_at", "updated_at"}).
		AddRow(movementID1, searchInvoiceID, 40.0, "DEBIT", "Search Invoice and Status", now, "PENDING", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements" WHERE invoice_id = $1 AND status = $2`)).
		WithArgs(searchInvoiceID, string(searchStatus)).
		WillReturnRows(rows)

	results, err := repo.Search(ctx, criteria)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	if len(results) == 1 {
		assert.Equal(t, movementID1, results[0].MovementID)
		assert.Equal(t, searchInvoiceID, results[0].InvoiceID)
		assert.Equal(t, model.StatusPending, results[0].Status)
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMovementSQLRepository_Search_NoCriteria(t *testing.T) {
	gormDB, mock := setupMockDB(t)
	logger := zerolog.Nop()
	repo := persistence.NewMovementSQLRepository(gormDB, logger)

	ctx := context.Background()
	criteria := &model.SearchCriteria{} // Empty criteria

	now := time.Now().UTC().Truncate(time.Microsecond)
	rows := sqlmock.NewRows([]string{"id", "invoice_id", "amount", "movement_type", "description", "transaction_date", "status", "created_at", "updated_at"}).
		AddRow(uuid.New(), uuid.New(), 50.0, "CREDIT", "Search All 1", now, "INVOICED", now, now).
		AddRow(uuid.New(), uuid.New(), 60.0, "DEBIT", "Search All 2", now, "PENDING", now, now)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT * FROM "movements"`)).
		WillReturnRows(rows)

	results, err := repo.Search(ctx, criteria)

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
