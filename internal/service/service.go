package service

import (
	"context"
	"finance-manager/internal/core"
	"finance-manager/internal/data"
	appErrors "finance-manager/internal/errors"
	"log/slog"
	"regexp"
)

type Service interface {
	GetUsers(ctx context.Context) ([]*core.User, error)
	GetTransactions(ctx context.Context) ([]*core.Transaction, error)
	CreateUser(ctx context.Context, name, email string) (*core.User, error)
}

type service struct {
	data   data.Data
	logger *slog.Logger
}

func NewService(data data.Data, logger *slog.Logger) Service {
	return &service{
		data:   data,
		logger: logger,
	}
}

func (s *service) GetUsers(ctx context.Context) ([]*core.User, error) {
	users, err := s.data.GetUsers(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch users", "error", err)
		return nil, err
	}
	return users, nil
}

func (s *service) GetTransactions(ctx context.Context) ([]*core.Transaction, error) {
	transactions, err := s.data.GetTransactions(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch transactions", "error", err)
		return nil, err
	}
	return transactions, nil
}

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

const (
	minNameLength = 2
	maxNameLength = 100
)

func (s *service) CreateUser(ctx context.Context, name, email string) (*core.User, error) {
	// Validate name
	if name == "" {
		return nil, appErrors.NewValidationError("name", "name is required")
	}
	if len(name) < minNameLength {
		return nil, appErrors.NewValidationError("name", "name must be at least 2 characters")
	}
	if len(name) > maxNameLength {
		return nil, appErrors.NewValidationError("name", "name must not exceed 100 characters")
	}

	// Validate email
	if email == "" {
		return nil, appErrors.NewValidationError("email", "email is required")
	}
	if !emailRegex.MatchString(email) {
		return nil, appErrors.NewValidationError("email", "invalid email format")
	}

	user, err := s.data.CreateUser(ctx, name, email)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, err
	}
	return user, nil
}
