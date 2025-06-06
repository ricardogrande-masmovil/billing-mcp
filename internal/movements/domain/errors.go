package domain

import "errors"

var (
	// ErrMovementNotFound is returned when a movement is not found.
	ErrMovementNotFound = errors.New("movement not found")
	// ErrInvalidMovementData is returned when movement data is invalid.
	ErrInvalidMovementData = errors.New("invalid movement data")
	// ErrMovementCreationFailed is returned when movement creation fails.
	ErrMovementCreationFailed = errors.New("movement creation failed")
	// ErrMovementUpdateFailed is returned when movement update fails.
	ErrMovementUpdateFailed = errors.New("movement update failed")
	// ErrMovementDeletionFailed is returned when movement deletion fails.
	ErrMovementDeletionFailed = errors.New("movement deletion failed")
)
