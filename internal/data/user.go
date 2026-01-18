package data

import (
	"context"
	"errors"
	"finance-manager/internal/core"
	appErrors "finance-manager/internal/errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserModel struct {
	db *pgxpool.Pool
}

func (u *UserModel) GetUsers(ctx context.Context) ([]*core.User, error) {
	q := `SELECT id, name, email, created_at FROM users`
	rows, err := u.db.Query(ctx, q)
	if err != nil {
		return nil, appErrors.WrapDatabaseError(err)
	}
	defer rows.Close()

	users := make([]*core.User, 0)

	for rows.Next() {
		user := &core.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			return nil, appErrors.WrapDatabaseError(err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, appErrors.WrapDatabaseError(err)
	}

	return users, nil
}

func (u *UserModel) CreateUser(ctx context.Context, name, email string) (*core.User, error) {
	q := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at`
	var user core.User
	if err := u.db.QueryRow(ctx, q, name, email).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, appErrors.ErrDuplicateEmail
		}
		return nil, appErrors.WrapDatabaseError(err)
	}
	return &user, nil
}
