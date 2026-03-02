package service

import (
	"context"
	"encoding/json"
	"finance-manager/internal/core"
	"finance-manager/internal/data"
	"log/slog"
	"time"
)

type Service interface {
	GetUsers(ctx context.Context) ([]*core.User, error)
	GetTransactions(ctx context.Context) ([]*core.Transaction, error)
	CreateUser(ctx context.Context, name, email string) (*core.User, error)
	CreateTransaction(ctx context.Context, userId, amount, categoryId int64, description string, metadata json.RawMessage, tags []string) (*core.Transaction, error)
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
	return s.data.GetUsers(ctx)
}

func (s *service) GetTransactions(ctx context.Context) ([]*core.Transaction, error) {
	return s.data.GetTransactions(ctx)
}

func (s *service) CreateUser(ctx context.Context, name, email string) (*core.User, error) {
	if err := validateCreateUser(name, email); err != nil {
		return nil, err
	}
	return s.data.CreateUser(ctx, name, email)
}

func (s *service) CreateTransaction(ctx context.Context, userId, amount, categoryId int64, description string, metadata json.RawMessage, tags []string) (*core.Transaction, error) {
	if err := validateCreateTransaction(userId, amount, categoryId, description, metadata, tags); err != nil {
		return nil, err
	}
	return s.data.CreateTransaction(ctx, userId, amount, categoryId, description, metadata, tags, time.Now())
}
