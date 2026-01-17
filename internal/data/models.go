package data

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Models struct {
	UserModel
	TransactionModel
}

func NewModels(db *pgxpool.Pool) *Models {
	return &Models{
		UserModel{db: db},
		TransactionModel{db: db},
	}
}
