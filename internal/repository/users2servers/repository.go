package users2servers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-bot/internal/db"
)

const table = "users2servers"

const (
	colUserID   = "user_id"
	colServerID = "server_id"
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
			colUserID,
			colServerID,
			colAddress,
		).
		Values(
			u2s.UserID,
			u2s.ServerID,
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

	row, err := r.db.DB().QueryContext(ctx, query, args...)
	defer func() {
		if clsErr := row.Close(); clsErr != nil {
			err = errors.Join(err, clsErr)
		}
	}()
	if err != nil {
		return fmt.Errorf("user_repository.create_user: %w", err)
	}

	return nil
}

func (r *Repository) GetAllAllocatedIPsByServerAlias(ctx context.Context, alias string) ([]string, error) {
	q, args, err := squirrel.Select(fmt.Sprintf("%s.%s", table, colAddress)).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Join("servers ON users2servers.server_id = servers.id").
		Where(squirrel.Eq{"servers.name": alias}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("users2servers_repository.get_all_allocated_ips_by_server_alias %w", err)
	}

	query := db.Query{
		Name:     "users2servers_repository.create_users2servers",
		QueryRaw: q,
	}

	rows, err := r.db.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("users2servers_repository.get_all_allocated_ips_by_server_alias %w", err)
	}
	defer func() {
		if clsErr := rows.Close(); clsErr != nil {
			err = errors.Join(err, clsErr)
		}
	}()

	var ips []string
	for rows.Next() {
		var ip string
		if err := rows.Scan(&ip); err != nil {
			return nil, fmt.Errorf("users2servers_repository.get_all_allocated_ips_by_server_alias %w", err)
		}

		ips = append(ips, ip)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("users2servers_repository.get_all_allocated_ips_by_server_alias: %w", err)
	}

	return ips, nil
}

func (r *Repository) GetFullInfo(ctx context.Context, userID int64) (*UsersServers, error) {
	q, args, err := squirrel.Select(
		"users.id as user_id",
		"users.username as username",
		"users.first_name as first_name",
		"users.last_name as last_name",
		"users.role as role",
		"users.public_key as user_public_key",
		"users.private_key as user_private_key",
		"users.state as state",
		"users2servers.address as user_address",
		"servers.id as server_id",
		"servers.name as server_name",
		"servers.address as server_address",
		"servers.public_key as server_public_key",
		"servers.private_key as server_private_key",
	).
		PlaceholderFormat(squirrel.Dollar).
		From(table).
		Join("servers ON users2servers.server_id = servers.id").
		Join("users ON users2servers.user_id = users.id").
		Where(squirrel.Eq{"users.id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("users2servers_repository.get_full_info %w", err)
	}

	query := db.Query{
		Name:     "users2servers_repository.get_full_info",
		QueryRaw: q,
	}

	var usersServers UsersServers
	if err := r.db.DB().GetContext(ctx, &usersServers, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			//nolint:nilnil
			return nil, nil
		}
		return nil, fmt.Errorf("users2servers_repository.get_full_info %w", err)
	}

	return &usersServers, nil
}
