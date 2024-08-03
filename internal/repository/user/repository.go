package user

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-api/internal/db"
)

const table = "users"

const (
	colPk         = "id"
	colUsername   = "username"
	colFirstname  = "first_name"
	colLastname   = "last_name"
	colRole       = "role"
	colPublicKey  = "public_key"
	colPrivateKey = "private_key"
	colCreatedAt  = "created_at"
	colUpdatedAt  = "updated_at"
)

type Repository struct {
	db db.Client
}

func NewRepository(db db.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetUserById(ctx context.Context, id int64) (*Model, error) {
	return nil, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *Model) error {
	sql, args, err := squirrel.Insert(table).
		PlaceholderFormat(squirrel.Dollar).
		Columns(
			colPk,
			colUsername,
			colFirstname,
			colLastname,
			colRole,
			colPublicKey,
			colPrivateKey,
		).
		Values(
			user.Id,
			user.Username,
			user.FirstName,
			user.LastName,
			user.Role,
			user.PublicKey,
			user.PrivateKey,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("user_repository.create_user: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.create_user",
		QueryRaw: sql,
	}

	_, err = r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("user_repository.create_user: %w", err)
	}

	return nil
}
