package sql

import (
	"github.com/ricardogrande-masmovil/billing-mcp/internal/invoices/domain/model"
)

type InvoiceSqlConverter struct {
}

func NewInvoiceSqlConverter() InvoiceSqlConverter {
	return InvoiceSqlConverter{}
}

func (c InvoiceSqlConverter) ConvertInvoiceToDomain(invoice Invoice) (model.Invoice, error) {
	domainStatus, err := model.GetStatusFromString(invoice.Status)
	if err != nil {
		return model.Invoice{}, err
	}

	return model.Invoice{
		ID:                    model.InvoiceID(invoice.ID),
		AccountID:             invoice.AccountID,
		IssueDate:             invoice.IssueDate,
		DueDate:               invoice.DueDate,
		TaxAmount:             invoice.TaxAmount,
		TotalAmountWithoutTax: invoice.TotalAmountWithoutTax,
		TotalAmountWithTax:    invoice.TotalAmountWithTax,
		Status:                domainStatus,
		InvoiceNumber:         invoice.InvoiceNumber,
	}, nil
}

func (c InvoiceSqlConverter) ConvertInvoicesToDomain(invoices []Invoice) ([]model.Invoice, error) {
	domainInvoices := make([]model.Invoice, len(invoices))
	for i, invoice := range invoices {
		domainInvoice, err := c.ConvertInvoiceToDomain(invoice)
		if err != nil {
			return nil, err
		}
		domainInvoices[i] = domainInvoice
	}
	return domainInvoices, nil
}

func (c InvoiceSqlConverter) ConvertCriteriaToSql(criteria model.Criteria) map[string]interface{} {
	sqlCriteria := make(map[string]interface{})

	if criteria.Status != "" {
		sqlCriteria["status"] = criteria.Status
	}

	if !criteria.IssueDateFrom.IsZero() {
		sqlCriteria["issue_date_from"] = criteria.IssueDateFrom
	}

	if !criteria.IssueDateTo.IsZero() {
		sqlCriteria["issue_date_to"] = criteria.IssueDateTo
	}

	return sqlCriteria
}

// SQLLineToInvoiceLine converts a SQL model InvoiceLine to a domain InvoiceLine
func (c InvoiceSqlConverter) SQLLineToInvoiceLine(line InvoiceLine) model.InvoiceLine {
	return model.InvoiceLine{
		MovementID:       line.MovementID,
		Description:      line.Description,
		AmountWithoutTax: line.AmountWithoutTax,
		AmountWithTax:    line.AmountWithTax,
		TaxPercentage:    line.TaxPercentage,
		OperationType:    line.OperationType,
	}
}
