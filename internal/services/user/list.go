package user

import (
	"context"
	"fmt"

	"wireguard-bot/internal/repository/user"
)

func (u *ServiceUser) List(ctx context.Context) ([]user.Model, error) {
	users, err := u.userRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("user_service.list %w", err)
	}

	return users, nil
}
