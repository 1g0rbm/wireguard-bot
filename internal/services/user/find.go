package user

import (
	"context"
	"fmt"

	"wireguard-api/internal/repository/user"
)

func (u *ServiceUser) FindUser(ctx context.Context, userID int64) (*user.Model, error) {
	model, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[user_service.is_user_exist] %w", err)
	}

	return model, nil
}
