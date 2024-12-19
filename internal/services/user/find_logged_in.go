package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"wireguard-bot/internal/repository/user"
)

func (u *ServiceUser) FindLoggedIn(ctx context.Context, sessionID uuid.UUID) (*user.Model, error) {
	model, err := u.userRepo.FIndUserByActiveSessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("[user_service.find_logged_in] %w", err)
	}

	return model, nil
}
