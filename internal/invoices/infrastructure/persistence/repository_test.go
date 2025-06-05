package persistence

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/infrastructure/persistence/sql"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockInvoiceSqlClient is a mock implementation of the InvoiceSqlClient
type MockInvoiceSqlClient struct {
	mock.Mock
}

func (m *MockInvoiceSqlClient) GetInvoiceByID(id string) (sql.InvoiceModel, error) {
	args := m.Called(id)
	return args.Get(0).(sql.InvoiceModel), args.Error(1)
}

func (m *MockInvoiceSqlClient) GetInvoicesByAccountId(accountId string, criteria sql.SearchCriteria) ([]sql.InvoiceModel, error) {
	args := m.Called(accountId, criteria)
	return args.Get(0).([]sql.InvoiceModel), args.Error(1)
}

func (m *MockInvoiceSqlClient) GetInvoiceLinesByInvoiceID(ctx context.Context, invoiceID string) ([]sql.InvoiceLineModel, error) {
	args := m.Called(ctx, invoiceID)
	return args.Get(0).([]sql.InvoiceLineModel), args.Error(1)
}

// MockInvoiceSqlConverter is a mock implementation of the InvoiceSqlConverter
type MockInvoiceSqlConverter struct {
	mock.Mock
}

func (m *MockInvoiceSqlConverter) ConvertInvoiceToDomain(sqlModel sql.InvoiceModel) (domain.Invoice, error) {
	args := m.Called(sqlModel)
	return args.Get(0).(domain.Invoice), args.Error(1)
}

func (m *MockInvoiceSqlConverter) ConvertInvoicesToDomain(sqlModels []sql.InvoiceModel) (domain.Invoices, error) {
	args := m.Called(sqlModels)
	return args.Get(0).(domain.Invoices), args.Error(1)
}

func (m *MockInvoiceSqlConverter) ConvertCriteriaToSql(domainCriteria domain.Criteria) sql.SearchCriteria {
	args := m.Called(domainCriteria)
	return args.Get(0).(sql.SearchCriteria)
}

func (m *MockInvoiceSqlConverter) SQLLineToInvoiceLine(sqlLine sql.InvoiceLineModel) domain.InvoiceLine {
	args := m.Called(sqlLine)
	return args.Get(0).(domain.InvoiceLine)
}

func TestGetInvoiceLines(t *testing.T) {
	// Setup
	mockClient := new(MockInvoiceSqlClient)
	mockConverter := new(MockInvoiceSqlConverter)
	logger := zerolog.New(zerolog.NewTestWriter(t))
	repository := Repository{
		invoiceSqlClient: mockClient,
		converter:        mockConverter,
		logger:           logger,
	}

	ctx := context.Background()
	invoiceID := domain.InvoiceID(uuid.New())

	t.Run("Successfully retrieves invoice lines", func(t *testing.T) {
		// Setup mock SQL invoice lines
		sqlLines := []sql.InvoiceLineModel{
			{
				ID:               uuid.New(),
				InvoiceID:        uuid.MustParse(invoiceID.String()),
				MovementID:       uuid.New(),
				Description:      "Test Line 1",
				AmountWithoutTax: 100.0,
				AmountWithTax:    121.0,
				TaxPercentage:    21.0,
				OperationType:    "CHARGE",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               uuid.New(),
				InvoiceID:        uuid.MustParse(invoiceID.String()),
				MovementID:       uuid.New(),
				Description:      "Test Line 2",
				AmountWithoutTax: 50.0,
				AmountWithTax:    60.5,
				TaxPercentage:    21.0,
				OperationType:    "DISCOUNT",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}

		// Setup mock domain invoice lines
		domainLines := []domain.InvoiceLine{
			{
				MovementID:       sqlLines[0].MovementID,
				Description:      sqlLines[0].Description,
				AmountWithoutTax: sqlLines[0].AmountWithoutTax,
				AmountWithTax:    sqlLines[0].AmountWithTax,
				TaxPercentage:    sqlLines[0].TaxPercentage,
				OperationType:    sqlLines[0].OperationType,
			},
			{
				MovementID:       sqlLines[1].MovementID,
				Description:      sqlLines[1].Description,
				AmountWithoutTax: sqlLines[1].AmountWithoutTax,
				AmountWithTax:    sqlLines[1].AmountWithTax,
				TaxPercentage:    sqlLines[1].TaxPercentage,
				OperationType:    sqlLines[1].OperationType,
			},
		}

		// Set expectations
		mockClient.On("GetInvoiceLinesByInvoiceID", ctx, invoiceID.String()).Return(sqlLines, nil)
		mockConverter.On("SQLLineToInvoiceLine", sqlLines[0]).Return(domainLines[0])
		mockConverter.On("SQLLineToInvoiceLine", sqlLines[1]).Return(domainLines[1])

		// Execute
		result, err := repository.GetInvoiceLines(ctx, invoiceID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, 2, len(result))
		assert.Equal(t, domainLines[0].MovementID, result[0].MovementID)
		assert.Equal(t, domainLines[0].Description, result[0].Description)
		assert.Equal(t, domainLines[0].AmountWithoutTax, result[0].AmountWithoutTax)
		assert.Equal(t, domainLines[0].AmountWithTax, result[0].AmountWithTax)
		assert.Equal(t, domainLines[0].TaxPercentage, result[0].TaxPercentage)
		assert.Equal(t, domainLines[0].OperationType, result[0].OperationType)

		// Verify second line
		assert.Equal(t, domainLines[1].MovementID, result[1].MovementID)
		assert.Equal(t, domainLines[1].Description, result[1].Description)

		// Verify expectations
		mockClient.AssertExpectations(t)
		mockConverter.AssertExpectations(t)
	})

	t.Run("Handles empty invoice lines", func(t *testing.T) {
		// Reset mocks
		mockClient = new(MockInvoiceSqlClient)
		mockConverter = new(MockInvoiceSqlConverter)
		repository = Repository{
			invoiceSqlClient: mockClient,
			converter:        mockConverter,
			logger:           logger,
		}

		// Setup
		emptyLines := []sql.InvoiceLineModel{}

		// Set expectations
		mockClient.On("GetInvoiceLinesByInvoiceID", ctx, invoiceID.String()).Return(emptyLines, nil)

		// Execute
		result, err := repository.GetInvoiceLines(ctx, invoiceID)

		// Assert
		require.NoError(t, err)
		assert.Empty(t, result)

		// Verify expectations
		mockClient.AssertExpectations(t)
	})
}
