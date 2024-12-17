package services

import (
	"context"

	"github.com/google/uuid"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/utils"
)

type ConfigService interface {
	GenerateConf(ctx context.Context, userID int64) ([]byte, error)
	GenerateQR(ctx context.Context, userID int64) ([]byte, error)
}

type UserService interface {
	Create(ctx context.Context, serverID int, userID int64, username, firstName, lastname string) (*user.Model, error)
	Find(ctx context.Context, userID int64) (*user.Model, error)
	GetOrCreate(ctx context.Context, serverID int, userID int64, username, firstName, lastname string) (*user.Model, error)
	LoginAdmin(ctx context.Context, username string) error
	List(ctx context.Context, opts ...utils.FilterOption) ([]user.Model, error)
	Enable(ctx context.Context, userID int64) error
}

type SessionService interface {
	CreateOrUpdate(ctx context.Context, userID int64) error
	Check(ctx context.Context, sessionID uuid.UUID) error
	CheckByUsername(ctx context.Context, username string) (*uuid.UUID, error)
}
