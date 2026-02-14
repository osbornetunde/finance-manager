package service

import (
	appErrors "finance-manager/internal/errors"
	"regexp"
	"strings"
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

func validateCreateTransaction(userID, amount, categoryID int64, description string) error {
	if userID <= 0 {
		return appErrors.NewValidationError("user_id", "user_id must be greater than 0")
	}
	if amount <= 0 {
		return appErrors.NewValidationError("amount", "amount must be a positive integer")
	}
	if categoryID <= 0 {
		return appErrors.NewValidationError("category_id", "category_id must be greater than 0")
	}
	if strings.TrimSpace(description) == "" {
		return appErrors.NewValidationError("description", "description is required")
	}
	return nil
}
