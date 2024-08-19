package services

import "context"

type ConfigService interface {
	GenerateConf(ctx context.Context, userID int64) ([]byte, error)
	GenerateQR(ctx context.Context, userID int64) ([]byte, error)
}

type UserService interface {
	Create(ctx context.Context, serverID int, userID int64, username, firstName, lastname string) error
	IsUserExist(ctx context.Context, userID int64) (bool, error)
}
