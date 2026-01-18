package errors

import (
	"errors"
	"fmt"
)

var (
	// User errors
	ErrDuplicateEmail = errors.New("email already exists")

	// Generic errors
	ErrDatabaseOperation = errors.New("database operation failed")
)

// WrapDatabaseError wraps a database error while preserving the original for logging
func WrapDatabaseError(err error) error {
	return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
}

// ValidationError represents a validation failure on a specific field
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

// NewValidationError creates a new validation error for a specific field
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

// IsDuplicateEmail checks if an error is a duplicate email error
func IsDuplicateEmail(err error) bool {
	return errors.Is(err, ErrDuplicateEmail)
}

// IsDatabaseError checks if an error is a database operation error
func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}
