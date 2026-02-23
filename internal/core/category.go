package core

import (
	"encoding/json"
	"fmt"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

// IsValid checks if the CategoryType is a valid value
func (c CategoryType) IsValid() bool {
	switch c {
	case CategoryTypeIncome, CategoryTypeExpense:
		return true
	}
	return false
}

// UnmarshalJSON validates the CategoryType when parsing JSON
func (c *CategoryType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	ct := CategoryType(s)
	if !ct.IsValid() {
		return fmt.Errorf("invalid category type: %q, must be %q or %q",
			s, CategoryTypeIncome, CategoryTypeExpense)
	}
	*c = ct
	return nil
}

type Category struct {
	ID     int64        `json:"id"`
	Name   string       `json:"name"`
	Type   CategoryType `json:"type"`
	UserID int64        `json:"user_id"`
}
