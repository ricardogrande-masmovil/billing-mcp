package ports

type Invoices = []Invoice

type Invoice struct {
	ID               string `json:"id"`
	AmountWithoutTax int    `json:"amount_without_tax"`
	AmountWithTax    int    `json:"amount_with_tax"`
	Status           string `json:"status"`
	IssueDate        string `json:"issue_date"`
	DueDate          string `json:"due_date"`
}

// InvoiceMovementDTO represents a simplified view of a movement in the context of an invoice
type InvoiceMovementDTO struct {
	MovementID       string  `json:"movement_id"`
	Description      string  `json:"description"`
	Amount           float64 `json:"amount"`
	AmountWithoutTax float64 `json:"amount_without_tax"`
	AmountWithTax    float64 `json:"amount_with_tax"`
	TaxPercentage    float64 `json:"tax_percentage"`
	OperationType    string  `json:"operation_type"`
}

// InvoiceMovementsDTO is a slice of InvoiceMovementDTO
type InvoiceMovementsDTO = []InvoiceMovementDTO
