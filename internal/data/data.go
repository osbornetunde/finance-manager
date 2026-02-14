package data

import (
	"context"
	"encoding/json"
	"finance-manager/internal/core"
)

type Data interface {
	GetUsers(ctx context.Context) ([]*core.User, error)
	GetTransactions(ctx context.Context) ([]*core.Transaction, error)
	CreateUser(ctx context.Context, name, email string) (*core.User, error)
	CreateTransaction(ctx context.Context, userId, amount, categoryId int64,
		description string, metadata json.RawMessage, tags []string) (*core.Transaction, error)
}
