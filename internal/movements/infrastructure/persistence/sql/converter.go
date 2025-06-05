package sql

import (
	domainmodel "github.com/ricardogrande-masmovil/billing-mcp/internal/movements/domain/model"
	"github.com/ricardogrande-masmovil/billing-mcp/pkg/persistence"
)

// MovementConverter handles mapping between domain and SQL movement models.
type MovementConverter struct{}

// NewMovementConverter creates a new MovementConverter.
func NewMovementConverter() *MovementConverter {
	return &MovementConverter{}
}

// ToDomainMovement converts an SQL movement model to a domain movement model.
func (c *MovementConverter) ToDomainMovement(sqlMovement *Movement) *domainmodel.Movement {
	if sqlMovement == nil {
		return nil
	}
	status, _ := domainmodel.StatusFromString(sqlMovement.Status)                   // Handle error appropriately
	movementType, _ := domainmodel.MovementTypeFromString(sqlMovement.MovementType) // Handle error appropriately

	return &domainmodel.Movement{
		MovementID:      sqlMovement.ID,
		InvoiceID:       sqlMovement.InvoiceID,
		Amount:          sqlMovement.Amount,
		MovementType:    movementType,
		Description:     sqlMovement.Description,
		TransactionDate: sqlMovement.TransactionDate,
		Status:          status,
	}
}

// ToSQLMovement converts a domain movement model to an SQL movement model.
func (c *MovementConverter) ToSQLMovement(domainMovement *domainmodel.Movement) *Movement {
	if domainMovement == nil {
		return nil
	}
	return &Movement{
		BaseModel: persistence.BaseModel{
			ID:        domainMovement.MovementID,
		},
		InvoiceID:       domainMovement.InvoiceID,
		Amount:          domainMovement.Amount,
		MovementType:    domainMovement.MovementType.String(),
		Description:     domainMovement.Description,
		TransactionDate: domainMovement.TransactionDate,
		Status:          domainMovement.Status.String(),
	}
}
