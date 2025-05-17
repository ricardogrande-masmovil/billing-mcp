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