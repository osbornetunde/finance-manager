package data

import (
	"context"
	"encoding/json"
	"errors"
	"finance-manager/internal/core"
	appErrors "finance-manager/internal/errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionModel struct {
	db *pgxpool.Pool
}

func (t *TransactionModel) GetTransactions(ctx context.Context) ([]*core.Transaction, error) {
	q := `SELECT id, user_id, amount, description, date, category_id, metadata, tags FROM transactions`
	rows, err := t.db.Query(ctx, q)
	if err != nil {
		return nil, appErrors.WrapDatabaseError(err)
	}
	defer rows.Close()

	var transactions []*core.Transaction
	for rows.Next() {
		var transaction core.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CategoryID, &transaction.Metadata, &transaction.Tags); err != nil {
			return nil, appErrors.WrapDatabaseError(err)
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, appErrors.WrapDatabaseError(err)
	}
	return transactions, nil
}

func (t *TransactionModel) CreateTransaction(ctx context.Context, userId, amount, categoryId int64, description string, metadata json.RawMessage, tags []string, createdAt time.Time) (*core.Transaction, error) {
	q := `INSERT INTO transactions (user_id, amount, description, date, category_id, metadata, tags) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, user_id, amount, description, date, category_id, metadata, tags`

	var transaction core.Transaction
	if err := t.db.QueryRow(ctx, q, userId, amount, description, createdAt,
		categoryId, metadata, tags).Scan(&transaction.ID,
		&transaction.UserID, &transaction.Amount,
		&transaction.Description, &transaction.Date, &transaction.CategoryID, &transaction.Metadata, &transaction.Tags); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			constraint := pgErr.ConstraintName
			if constraint == "transactions_user_id_fkey" {
				return nil, appErrors.ErrInvalidUserReference
			}
			if constraint == "transactions_category_id_fkey" {
				return nil, appErrors.ErrInvalidCategoryReference
			}
		}
		return nil, appErrors.WrapDatabaseError(err)
	}
	return &transaction, nil
}
