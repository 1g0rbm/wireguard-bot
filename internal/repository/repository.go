package repository

import (
	"context"

	"wireguard-api/internal/repository/user"
)

type UserRepository interface {
	GetUserById(ctx context.Context, id int64) (*user.Model, error)
	CreateUser(ctx context.Context, user *user.Model) error
}
