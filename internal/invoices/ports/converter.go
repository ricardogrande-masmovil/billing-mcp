package ports

import (
	"encoding/json"
	"errors"

	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
)

var (
	ErrInvalidInvoice      = errors.New("invalid invoice")
	ErrInvalidInvoices     = errors.New("invalid invoices")
)

type Converter struct{}

func NewConverter() Converter {
	return Converter{}
}

func (c Converter) ConvertDomainInvoiceToJsonInvoice(domainInvoice domain.Invoice) ([]byte, error) {
	mcpInvoice := Invoice{
		ID:               string(domainInvoice.ID),
		AmountWithoutTax: int(domainInvoice.TotalAmountWithoutTax),
		AmountWithTax:    int(domainInvoice.TotalAmountWithTax),
		Status:           string(domainInvoice.Status),
		IssueDate:        domainInvoice.IssueDate.Format("2006-01-02"),
		DueDate:          domainInvoice.DueDate.Format("2006-01-02"),
	}
	jsonData, err := json.Marshal(mcpInvoice)
	if err != nil {
		return nil, ErrInvalidInvoice
	}
	return jsonData, nil
}

func (c Converter) ConvertDomainInvoicesToJsonInvoices(domainInvoices domain.Invoices) ([]byte, error) {
	mcpInvoices := make([]Invoice, len(domainInvoices))
	for i, domainInvoice := range domainInvoices {
		mcpInvoices[i] = Invoice{
			ID:               string(domainInvoice.ID),
			AmountWithoutTax: int(domainInvoice.TotalAmountWithoutTax),
			AmountWithTax:    int(domainInvoice.TotalAmountWithTax),
			Status:           string(domainInvoice.Status),
			IssueDate:        domainInvoice.IssueDate.Format("2006-01-02"),
			DueDate:          domainInvoice.DueDate.Format("2006-01-02"),
		}
	}
	jsonData, err := json.Marshal(mcpInvoices)
	if err != nil {
		return nil, ErrInvalidInvoices
	}
	return jsonData, nil
}