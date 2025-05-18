package model

import (
	"errors"
)

var ErrStatusUnknown = errors.New("unknown status")

// InvoiceStatus represents the status of an invoice.
type InvoiceStatus string

const (
	InvoiceStatusDraft   InvoiceStatus = "DRAFT"
	InvoiceStatusSent    InvoiceStatus = "SENT"
	InvoiceStatusPaid    InvoiceStatus = "PAID"
	InvoiceStatusOverdue InvoiceStatus = "OVERDUE"
	InvoiceStatusVoid    InvoiceStatus = "VOID"
	InvoiceStatusUnpaid  InvoiceStatus = "UNPAID"
)

var statusStringMap = map[string]InvoiceStatus{
	"DRAFT":   InvoiceStatusDraft,
	"SENT":    InvoiceStatusSent,
	"PAID":    InvoiceStatusPaid,
	"OVERDUE": InvoiceStatusOverdue,
	"VOID":    InvoiceStatusVoid,
	"UNPAID":  InvoiceStatusUnpaid,
}

func GetStatusFromString(status string) (InvoiceStatus, error) {
	if s, ok := statusStringMap[status]; ok {
		return s, nil
	}
	return "", ErrStatusUnknown
}