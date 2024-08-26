package services

import (
	"context"

	"wireguard-api/internal/repository/user"
)

type ConfigService interface {
	GenerateConf(ctx context.Context, userID int64) ([]byte, error)
	GenerateQR(ctx context.Context, userID int64) ([]byte, error)
}

type UserService interface {
	Create(ctx context.Context, serverID int, userID int64, username, firstName, lastname string) error
	FindUser(ctx context.Context, userID int64) (*user.Model, error)
}
