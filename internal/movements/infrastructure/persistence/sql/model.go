package sql

import (
	"time"

	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/pkg/persistence"
)

// Movement is the GORM model for an invoice movement.
// It includes database-specific fields like timestamps and soft delete.
// It maps to the "movements" table in the database.
type Movement struct {
	persistence.BaseModel
	InvoiceID       uuid.UUID `gorm:"type:uuid;not null"`
	Amount          float64   `gorm:"type:decimal(10,2);not null"`
	MovementType    string    `gorm:"type:varchar(50);not null"`
	Description     string    `gorm:"type:text"`
	TransactionDate time.Time `gorm:"not null"`
	Status          string    `gorm:"type:varchar(50);not null"`
}

// TableName specifies the table name for the Movement model.
func (Movement) TableName() string {
	return "movements"
}
