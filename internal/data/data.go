package data

import (
	"context"
	"finance-manager/internal/core"
)

type Data interface {
	GetUsers(ctx context.Context) ([]*core.User, error)
	GetTransactions(ctx context.Context) ([]*core.Transaction, error)
	CreateUser(ctx context.Context, name, email string) (*core.User, error)
}
