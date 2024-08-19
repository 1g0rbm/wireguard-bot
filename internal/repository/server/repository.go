package server

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-api/internal/db"
)

const table = "servers"

const (
	colPk         = "id"
	colName       = "name"
	colAddress    = "address"
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

func (r *Repository) GetByName(ctx context.Context, name string) (*Server, error) {
	q, args, err := squirrel.Select(
		colPk,
		colName,
		colAddress,
		colPublicKey,
		colPrivateKey,
		colCreatedAt,
		colUpdatedAt,
	).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Where(squirrel.Eq{colName: name}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("server_repository.get_server_by_name: %w", err)
	}

	query := db.Query{
		Name:     "server_repository.Get",
		QueryRaw: q,
	}

	var model Server
	if err = r.db.DB().GetContext(ctx, &model, query, args...); err != nil {
		return nil, fmt.Errorf("server_repository.get_server_by_name: %w", err)
	}

	return &model, nil
}
