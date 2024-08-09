package services

import "context"

type ConfigService interface {
	GenerateConf(ctx context.Context, userId int64) ([]byte, error)
}
