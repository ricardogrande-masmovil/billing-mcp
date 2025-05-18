package ports

import (
	"encoding/json"
	"errors"
	"time"

	domain "github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
)

var (
	ErrInvalidInvoice      = errors.New("invalid invoice")
	ErrInvalidInvoices     = errors.New("invalid invoices")
	ErrInvalidStatusCriteria = errors.New("invalid status criteria")
	ErrInvalidDateCriteria = errors.New("invalid date criteria")
)

type Converter struct{}

func NewConverter() Converter {
	return Converter{}
}

func (c Converter) ConvertDomainInvoiceToJsonInvoice(domainInvoice domain.Invoice) ([]byte, error) {
	mcpInvoice := Invoice{
		ID:               domainInvoice.ID.String(),
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
			ID:               domainInvoice.ID.String(),
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

func (c Converter) ConvertRequestArgsToCriteria(args map[string]any) (domain.Criteria, error) {
	criteria := domain.Criteria{}

	paramStatus, ok := (args["status"]).(string)
	if !ok {
		return criteria, ErrInvalidStatusCriteria
	}

	if status, err := domain.GetStatusFromString(paramStatus); err != nil {
		return criteria, err
	} else {
		criteria.Status = status
	}

	paramStartDate, ok := (args["issueDateFrom"]).(string)
	if !ok {
		return criteria, ErrInvalidDateCriteria
	}
	if startDate, err := time.Parse(time.RFC3339, paramStartDate); err != nil {
		return criteria, err
	} else {
		criteria.IssueDateFrom = startDate
	}

	paramEndDate, ok := (args["issueDateTo"]).(string)
	if !ok {
		return criteria, ErrInvalidDateCriteria
	}
	if endDate, err := time.Parse(time.RFC3339, paramEndDate); err != nil {
		return criteria, err
	} else {
		criteria.IssueDateTo = endDate
	}

	return criteria, nil
}