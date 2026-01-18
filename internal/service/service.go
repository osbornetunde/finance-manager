package service

import (
	"context"
	"finance-manager/internal/core"
	"finance-manager/internal/data"
	"log/slog"
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

func (s *service) CreateUser(ctx context.Context, name, email string) (*core.User, error) {
	// Validate input
	if err := validateCreateUser(name, email); err != nil {
		return nil, err
	}

	user, err := s.data.CreateUser(ctx, name, email)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return nil, err
	}
	return user, nil
}
