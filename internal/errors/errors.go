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

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func IsValidationError(err error) bool {
	var validationErr *ValidationError
	return errors.As(err, &validationErr)
}

func IsDuplicateEmail(err error) bool {
	return errors.Is(err, ErrDuplicateEmail)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}
