package user

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"

	"wireguard-bot/internal/repository/user"
	"wireguard-bot/internal/utils"
)

func WithEnabled() utils.FilterOption {
	return func(f utils.Filter) {
		f.Eq["state"] = user.EnabledState
	}
}

func WithDisabled() utils.FilterOption {
	return func(f utils.Filter) {
		f.Eq["state"] = user.DisabledState
	}
}

func WithUsername(username string) utils.FilterOption {
	return func(f utils.Filter) {
		f.Like["username"] = username + "%"
	}
}

func (u *ServiceUser) List(ctx context.Context, opts ...utils.FilterOption) ([]user.Model, error) {
	f := utils.Filter{
		Eq:   squirrel.Eq{},
		Like: squirrel.Like{},
	}
	for _, opt := range opts {
		opt(f)
	}

	users, err := u.userRepo.List(ctx, f)
	if err != nil {
		return nil, fmt.Errorf("user_service.list %w", err)
	}

	return users, nil
}
