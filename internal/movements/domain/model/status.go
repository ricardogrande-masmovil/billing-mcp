package model

// Status represents the status of a movement.
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusInvoiced  Status = "INVOICED"
	StatusCancelled Status = "CANCELLED"
)
