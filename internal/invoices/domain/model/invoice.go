package domain

import (
	"errors"
	"time"
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
)

// InvoiceID represents the unique identifier for an Invoice.
type InvoiceID string

// InvoiceStatus represents the status of an invoice.
type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "Draft"
	InvoiceStatusSent    InvoiceStatus = "Sent"
	InvoiceStatusPaid    InvoiceStatus = "Paid"
	InvoiceStatusOverdue InvoiceStatus = "Overdue"
	InvoiceStatusVoid    InvoiceStatus = "Void"
	InvoiceStatusUnpaid  InvoiceStatus = "Unpaid"
)

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
