package model

import (
	"fmt"
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
	CreatedAt       time.Time // Added CreatedAt
	UpdatedAt       time.Time // Added UpdatedAt
}

// MovementType defines the type of movement (credit or debit).
type MovementType string

const (
	MovementTypeCredit MovementType = "CREDIT"
	MovementTypeDebit  MovementType = "DEBIT"
)

// String returns the string representation of the MovementType.
func (mt MovementType) String() string {
	return string(mt)
}

// MovementTypeFromString converts a string to a MovementType.
// Returns an error if the string is not a valid MovementType.
func MovementTypeFromString(s string) (MovementType, error) {
	switch s {
	case string(MovementTypeCredit):
		return MovementTypeCredit, nil
	case string(MovementTypeDebit):
		return MovementTypeDebit, nil
	default:
		return "", fmt.Errorf("invalid movement type: %s", s)
	}
}

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
		CreatedAt:       time.Now(),    // Initialize CreatedAt
		UpdatedAt:       time.Now(),    // Initialize UpdatedAt
	}, nil
}
