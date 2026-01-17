package data

import (
	"context"
	"finance-manager/internal/core"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionModel struct {
	db *pgxpool.Pool
}

func (t *TransactionModel) GetTransactions(ctx context.Context) ([]*core.Transaction, error) {
	q := `SELECT id, user_id, amount, description, date, category_id, metadata, tags FROM transactions`
	rows, err := t.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*core.Transaction
	for rows.Next() {
		var transaction core.Transaction
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.Amount, &transaction.Description, &transaction.Date, &transaction.CategoryID, &transaction.Metadata, &transaction.Tags); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return transactions, nil
}
