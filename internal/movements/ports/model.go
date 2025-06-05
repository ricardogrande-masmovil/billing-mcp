package ports

// MovementDTO represents a movement as returned by the MCP API
type MovementDTO struct {
	ID              string  `json:"id"`
	InvoiceID       string  `json:"invoice_id"`
	Amount          float64 `json:"amount"`
	MovementType    string  `json:"movement_type"`
	Description     string  `json:"description"`
	TransactionDate string  `json:"transaction_date"`
	Status          string  `json:"status"`
}

// MovementsDTO is a slice of MovementDTO
type MovementsDTO = []MovementDTO

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
