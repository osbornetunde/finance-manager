package data

import (
	"context"
	"finance-manager/internal/core"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	db *pgxpool.Pool
}

func (u *UserModel) GetUsers(ctx context.Context) ([]*core.User, error) {
	q := `SELECT id, name, email, created_at FROM users`
	rows, err := u.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*core.User

	for rows.Next() {
		var user core.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
