package session

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserDoesNotHaveSession = errors.New("user does not have session")
	ErrExpiredSession         = errors.New("session expired")
)

func (s *ServiceSession) Check(ctx context.Context, userID int64) error {
	session, err := s.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("session_service.check.find_by_user_id %w", err)
	}

	if session == nil {
		return ErrUserDoesNotHaveSession
	}

	if session.ExpiredAt.Before(time.Now()) {
		return ErrExpiredSession
	}

	return nil
}
