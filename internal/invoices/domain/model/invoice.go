package model

import (
	"errors"
	"time"

	"github.com/google/uuid" // Import for UUID
)

// Predefined domain errors
var (
	ErrInvoiceIDEmpty            = errors.New("invoice ID cannot be empty")
	ErrAccountIDEmpty            = errors.New("account ID cannot be empty")
	ErrQuantityNotPositive       = errors.New("quantity must be positive")
	ErrInvoiceAlreadyPaid        = errors.New("invoice is already paid")
	ErrInvoiceNotDraft           = errors.New("invoice can only be marked as sent from draft status")
	ErrVoidInvoiceCannotBePaid   = errors.New("void invoice cannot be marked as paid")
	ErrPaidInvoiceCannotBeVoided = errors.New("paid invoice cannot be voided")
	ErrInvoiceNotFound           = errors.New("invoice not found") // Added
)

// InvoiceID represents the unique identifier for an Invoice.
type InvoiceID uuid.UUID

// NewInvoiceID generates a new unique InvoiceID.
func NewInvoiceID() InvoiceID {
	return InvoiceID(uuid.New())
}

// String returns the string representation of the InvoiceID.
func (id InvoiceID) String() string {
	return uuid.UUID(id).String()
}

// IsNil checks if the InvoiceID is the nil UUID.
func (id InvoiceID) IsNil() bool {
	return uuid.UUID(id) == uuid.Nil
}

// ParseInvoiceID parses a string into an InvoiceID.
// Returns uuid.Nil and error if parsing fails.
func ParseInvoiceID(s string) (InvoiceID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return InvoiceID(uuid.Nil), err
	}
	return InvoiceID(id), nil
}

// InvoiceLine represents a single line item on an invoice.
type InvoiceLine struct {
	Description      string
	AmountWithoutTax float64
	AmountWithTax    float64
	TaxPercentage    float64
	OperationType    string // Assuming this is a string for simplicity ("credit" or "debit")
}

type Invoices = []Invoice

// Invoice represents the aggregate root for an invoice.
type Invoice struct {
	ID                    InvoiceID
	AccountID             string
	IssueDate             time.Time
	DueDate               time.Time
	Lines                 []InvoiceLine
	TaxAmount             float64
	TotalAmountWithoutTax float64
	TotalAmountWithTax    float64
	Status                InvoiceStatus
	InvoiceNumber         string
}

// AddLine adds a new line item to the invoice.
func (inv *Invoice) AddLine(invoiceLine InvoiceLine) error {
	if inv.Status != InvoiceStatusDraft {
		return ErrInvoiceAlreadyPaid
	}

	inv.Lines = append(inv.Lines, invoiceLine)
	return nil
}

// recalculateTotals updates the subtotal, tax, and total amounts.
func (inv *Invoice) recalculateTotals() {
	inv.TaxAmount = 0
	inv.TotalAmountWithoutTax = 0
	inv.TotalAmountWithTax = 0

	for _, line := range inv.Lines {
		if line.OperationType == "credit" {
			inv.TotalAmountWithoutTax -= line.AmountWithoutTax
			inv.TotalAmountWithTax -= line.AmountWithTax
			inv.TaxAmount -= (line.AmountWithTax - line.AmountWithoutTax)
		} else if line.OperationType == "debit" {
			inv.TotalAmountWithoutTax += line.AmountWithoutTax
			inv.TotalAmountWithTax += line.AmountWithTax
			inv.TaxAmount += (line.AmountWithTax - line.AmountWithoutTax)
		}
	}
}

func (inv *Invoice) MarkAsSent() error {
	if inv.Status != InvoiceStatusDraft {
		return ErrInvoiceNotDraft
	}

	inv.Status = InvoiceStatusSent
	return nil
}

func (inv *Invoice) MarkAsPaid() error {
	if inv.Status == InvoiceStatusVoid {
		return ErrVoidInvoiceCannotBePaid
	}

	inv.Status = InvoiceStatusPaid
	return nil
}

func (inv *Invoice) MarkAsVoid() error {
	if inv.Status == InvoiceStatusPaid {
		return ErrPaidInvoiceCannotBeVoided
	}

	inv.Status = InvoiceStatusVoid
	return nil
}
