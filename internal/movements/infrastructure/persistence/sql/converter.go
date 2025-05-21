package sql

import (
	"github.com/google/uuid"
	"github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
)

// ToDomainMovement converts a GORM Movement model to a domain Movement model.
func ToDomainMovement(sqlMovement *Movement) *model.Movement {
	if sqlMovement == nil {
		return nil
	}
	return &model.Movement{
		MovementID:      sqlMovement.ID, // Assuming BaseModel's ID is the MovementID
		InvoiceID:       sqlMovement.InvoiceID,
		Amount:          sqlMovement.Amount,
		MovementType:    sqlMovement.MovementType,
		Description:     sqlMovement.Description,
		TransactionDate: sqlMovement.TransactionDate,
		Status:          sqlMovement.Status,
	}
}

// ToSQLMovement converts a domain Movement model to a GORM Movement model.
func ToSQLMovement(domainMovement *model.Movement) *Movement {
	if domainMovement == nil {
		return nil
	}
	sqlMovement := &Movement{
		InvoiceID:       domainMovement.InvoiceID,
		Amount:          domainMovement.Amount,
		MovementType:    domainMovement.MovementType,
		Description:     domainMovement.Description,
		TransactionDate: domainMovement.TransactionDate,
		Status:          domainMovement.Status,
	}
	// If the domain model has an ID, set it in the SQL model's BaseModel
	if domainMovement.MovementID != uuid.Nil {
		sqlMovement.ID = domainMovement.MovementID
	}
	return sqlMovement
}

// ToDomainMovements converts a slice of GORM Movement models to a slice of domain Movement models.
func ToDomainMovements(sqlMovements []*Movement) []*model.Movement {
	if sqlMovements == nil {
		return nil
	}
	domainMovements := make([]*model.Movement, len(sqlMovements))
	for i, m := range sqlMovements {
		domainMovements[i] = ToDomainMovement(m)
	}
	return domainMovements
}
