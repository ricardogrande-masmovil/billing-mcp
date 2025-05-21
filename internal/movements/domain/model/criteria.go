package model

import "github.com/google/uuid"

// SearchCriteria represents the criteria for searching movements.
type SearchCriteria struct {
	InvoiceID *uuid.UUID
	Status    *Status
	// Add other filter fields as needed
}
