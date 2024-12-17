package user

import (
	"context"
	"fmt"

	"wireguard-bot/internal/repository/user"
)

func (u *ServiceUser) GetOrCreate(
	ctx context.Context,
	serverID int,
	userID int64,
	username,
	firstName,
	lastname string,
) (*user.Model, error) {
	userModel, err := u.Find(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("service_user.get_or_create.find %w", err)
	}
	if userModel != nil {
		return userModel, nil
	}

	userModel, err = u.Create(ctx, serverID, userID, username, firstName, lastname)
	if err != nil {
		return nil, fmt.Errorf("service_user.get_or_create.create %w", err)
	}
	if userModel == nil {
		return nil, fmt.Errorf("user create repository return nil model")
	}

	return userModel, nil
}
