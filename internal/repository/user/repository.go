package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-bot/internal/db"
)

const table = "users"

const (
	colPk         = "id"
	colUsername   = "username"
	colFirstname  = "first_name"
	colLastname   = "last_name"
	colRole       = "role"
	colState      = "state"
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

func (r *Repository) GetUserByID(ctx context.Context, id int64) (*Model, error) {
	return r.getUserBy(ctx, squirrel.Eq{colPk: id})
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*Model, error) {
	return r.getUserBy(ctx, squirrel.Eq{colUsername: username})
}

func (r *Repository) getUserBy(ctx context.Context, eq squirrel.Eq) (*Model, error) {
	q, args, err := squirrel.Select(
		colPk,
		colUsername,
		colFirstname,
		colLastname,
		colRole,
		colState,
		colPublicKey,
		colPrivateKey,
		colCreatedAt,
		colUpdatedAt,
	).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Where(eq).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("user_repository.get_user_by: %w", err)
	}

	query := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: q,
	}

	var model Model
	if err = r.db.DB().GetContext(ctx, &model, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil
			return nil, nil
		}
		return nil, fmt.Errorf("user_repository.get_user_by_username: %w", err)
	}

	return &model, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *Model) error {
	q, args, err := squirrel.Insert(table).
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
			user.ID,
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
		QueryRaw: q,
	}

	row, err := r.db.DB().QueryContext(ctx, query, args...)
	defer func() {
		clsErr := row.Close()
		if clsErr != nil {
			err = errors.Join(err, clsErr)
		}
	}()

	if err != nil {
		return fmt.Errorf("user_repository.create_user: %w", err)
	}

	return nil
}
