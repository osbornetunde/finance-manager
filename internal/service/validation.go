package service

import (
	appErrors "finance-manager/internal/errors"
	"regexp"
)

const (
	minNameLength = 2
	maxNameLength = 100
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func validateName(name string) error {
	if name == "" {
		return appErrors.NewValidationError("name", "name is required")
	}
	if len(name) < minNameLength {
		return appErrors.NewValidationError("name", "name must be at least 2 characters")
	}
	if len(name) > maxNameLength {
		return appErrors.NewValidationError("name", "name must not exceed 100 characters")
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return appErrors.NewValidationError("email", "email is required")
	}
	if !emailRegex.MatchString(email) {
		return appErrors.NewValidationError("email", "invalid email format")
	}
	return nil
}

func validateCreateUser(name, email string) error {
	if err := validateName(name); err != nil {
		return err
	}
	if err := validateEmail(email); err != nil {
		return err
	}
	return nil
}
