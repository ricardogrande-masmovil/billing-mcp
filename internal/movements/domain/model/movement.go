package model

import (
	"time"

	"github.com/google/uuid"
)

// Movement represents an invoice movement.
type Movement struct {
	MovementID      uuid.UUID
	InvoiceID       uuid.UUID
	Amount          float64
	MovementType    MovementType
	Description     string
	TransactionDate time.Time
	Status          Status
}

// MovementType defines the type of movement (credit or debit).
type MovementType string

const (
	MovementTypeCredit MovementType = "CREDIT"
	MovementTypeDebit  MovementType = "DEBIT"
)

// NewMovement creates a new movement.
func NewMovement(invoiceID uuid.UUID, amount float64, movementType MovementType, description string) (*Movement, error) {
	// TODO: Add validation logic
	return &Movement{
		MovementID:      uuid.New(),
		InvoiceID:       invoiceID,
		Amount:          amount,
		MovementType:    movementType,
		Description:     description,
		TransactionDate: time.Now(),
		Status:          StatusPending, // Default status
	}, nil
}
