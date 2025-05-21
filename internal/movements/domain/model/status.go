package model

import "fmt"

// Status represents the status of a movement.
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusInvoiced  Status = "INVOICED"
	StatusCancelled Status = "CANCELLED"
)

// String returns the string representation of the Status.
func (s Status) String() string {
	return string(s)
}

// StatusFromString converts a string to a Status type.
// Returns an error if the string is not a valid Status.
func StatusFromString(s string) (Status, error) {
	switch s {
	case string(StatusPending):
		return StatusPending, nil
	case string(StatusInvoiced):
		return StatusInvoiced, nil
	case string(StatusCancelled):
		return StatusCancelled, nil
	default:
		return "", fmt.Errorf("invalid movement status: %s", s)
	}
}
