package sql

import (
	"time"

	commons "github.com/ricardogrande-masmovil/billing-mcp/internal/commons/persistence/gorm"
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
