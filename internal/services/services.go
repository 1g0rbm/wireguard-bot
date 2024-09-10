package services

import (
	"context"

	"wireguard-bot/internal/repository/user"
)

type ConfigService interface {
	GenerateConf(ctx context.Context, userID int64) ([]byte, error)
	GenerateQR(ctx context.Context, userID int64) ([]byte, error)
}

type UserService interface {
	Create(ctx context.Context, serverID int, userID int64, username, firstName, lastname string) error
	FindUser(ctx context.Context, userID int64) (*user.Model, error)
	LoginAdmin(ctx context.Context, username string) error
}

type SessionService interface {
	CreateOrUpdate(ctx context.Context, userID int64) error
	Check(ctx context.Context, userID int64) error
}
