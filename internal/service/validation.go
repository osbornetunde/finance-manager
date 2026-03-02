package service

import (
	"encoding/json"
	appErrors "finance-manager/internal/errors"
	"regexp"
	"strings"
)

const (
	minNameLength        = 2
	maxNameLength        = 100
	maxDescriptionLength = 500
	maxMetadataBytes     = 64 * 1024
	maxMetadataDepth     = 8
	maxTagCount          = 20
	maxTagLength         = 50
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

func validateCreateTransaction(userID, amount, categoryID int64, description string, metadata json.RawMessage, tags []string) error {
	if userID <= 0 {
		return appErrors.NewValidationError("user_id", "user_id must be greater than 0")
	}
	if amount <= 0 {
		return appErrors.NewValidationError("amount", "amount must be a positive integer")
	}
	if categoryID <= 0 {
		return appErrors.NewValidationError("category_id", "category_id must be greater than 0")
	}
	desc := strings.TrimSpace(description)
	if desc == "" {
		return appErrors.NewValidationError("description", "description is required")
	}
	if len(desc) > maxDescriptionLength {
		return appErrors.NewValidationError("description", "description must not exceed 500 characters")
	}
	if len(metadata) > maxMetadataBytes {
		return appErrors.NewValidationError("metadata", "metadata must not exceed 65536 bytes")
	}
	if len(metadata) > 0 {
		var parsed any
		if err := json.Unmarshal(metadata, &parsed); err != nil {
			return appErrors.NewValidationError("metadata", "metadata must be valid JSON")
		}
		if jsonDepth(parsed) > maxMetadataDepth {
			return appErrors.NewValidationError("metadata", "metadata nesting is too deep")
		}
	}
	if len(tags) > maxTagCount {
		return appErrors.NewValidationError("tags", "tags must not contain more than 20 values")
	}
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed == "" {
			return appErrors.NewValidationError("tags", "tags must not contain empty values")
		}
		if len(trimmed) > maxTagLength {
			return appErrors.NewValidationError("tags", "each tag must not exceed 50 characters")
		}
	}
	return nil
}

func jsonDepth(v any) int {
	switch val := v.(type) {
	case map[string]any:
		max := 1
		for _, child := range val {
			depth := 1 + jsonDepth(child)
			if depth > max {
				max = depth
			}
		}
		return max
	case []any:
		max := 1
		for _, child := range val {
			depth := 1 + jsonDepth(child)
			if depth > max {
				max = depth
			}
		}
		return max
	default:
		return 1
	}
}
