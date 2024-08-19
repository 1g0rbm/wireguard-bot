package users2servers

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-api/internal/db"
)

const table = "users2servers"

const (
	colUserId   = "user_id"
	colServerId = "server_id"
	colAddress  = "address"
)

type Repository struct {
	db db.Client
}

func NewRepository(db db.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUsers2Servers(ctx context.Context, u2s *Users2Servers) error {
	q, args, err := squirrel.Insert(table).
		PlaceholderFormat(squirrel.Dollar).
		Columns(
			colUserId,
			colServerId,
			colAddress,
		).
		Values(
			u2s.UserId,
			u2s.ServerId,
			u2s.Address,
		).
		ToSql()

	if err != nil {
		return fmt.Errorf("users2servers_repository.create_record: %w", err)
	}

	query := db.Query{
		Name:     "users2servers_repository.create_users2servers",
		QueryRaw: q,
	}

	_, err = r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("user_repository.create_user: %w", err)
	}

	return nil
}
