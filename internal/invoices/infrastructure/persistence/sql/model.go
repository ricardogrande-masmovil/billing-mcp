package sql

import (
	"time"

	"github.com/google/uuid"
	commons "github.com/ricardogrande-masmovil/billing-mcp/pkg/persistence"
)

// DBInvoice represents the invoice entity in the database.
// It includes GORM model fields (ID, CreatedAt, UpdatedAt, DeletedAt).
type Invoice struct {
	commons.BaseModel
	AccountID             string `gorm:"index"`
	IssueDate             time.Time
	DueDate               time.Time
	TaxAmount             float64
	TotalAmountWithoutTax float64
	TotalAmountWithTax    float64
	Status                string // e.g., "Draft", "Sent", "Paid", "Void"
	InvoiceNumber         string `gorm:"index;unique"`
}

// TableName specifies the table name for DBInvoice in the database.
func (Invoice) TableName() string {
	return "invoices"
}

// InvoiceLine represents a line item in an invoice, which corresponds to a movement
type InvoiceLine struct {
	MovementID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	InvoiceID        uuid.UUID `gorm:"type:uuid;not null;index"`
	Description      string    `gorm:"type:text"`
	AmountWithoutTax float64   `gorm:"type:decimal(10,2);not null"`
	AmountWithTax    float64   `gorm:"type:decimal(10,2);not null"`
	TaxPercentage    float64   `gorm:"type:decimal(5,2);not null"`
	OperationType    string    `gorm:"type:varchar(50);not null"` // "CREDIT" or "DEBIT"
}

// TableName specifies the table name for InvoiceLine in the database.
func (InvoiceLine) TableName() string {
	return "movements" // We're using the same table as movements since they represent the same data
}
