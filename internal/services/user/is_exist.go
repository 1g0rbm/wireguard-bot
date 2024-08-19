package user

import (
	"context"
	"fmt"
)

func (u *ServiceUser) IsUserExist(ctx context.Context, userId int64) (bool, error) {
	model, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return false, fmt.Errorf("[user_service.is_user_exist] %w", err)
	}

	return model != nil, nil
}
