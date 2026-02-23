package errors

import (
	"errors"
	"fmt"
)

var (
	// User errors
	ErrDuplicateEmail           = errors.New("email already exists")
	ErrInvalidUserReference     = errors.New("invalid user reference")
	ErrInvalidCategoryReference = errors.New("invalid category reference")

	// Generic errors
	ErrDatabaseOperation = errors.New("database operation failed")
)

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

func IsInvalidUserReference(err error) bool {
	return errors.Is(err, ErrInvalidUserReference)
}

func IsInvalidCategoryReference(err error) bool {
	return errors.Is(err, ErrInvalidCategoryReference)
}

func IsDatabaseError(err error) bool {
	return errors.Is(err, ErrDatabaseOperation)
}
