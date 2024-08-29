package session

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"wireguard-bot/internal/repository/session"
)

const sessionTTL = 24 * time.Hour

func (s *ServiceSession) CreateOrUpdate(ctx context.Context, userID int64) error {
	err := s.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		userSession, err := s.sessionRepo.FindByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("session_service.create_or_update %w", err)
		}

		if userSession == nil {
			userSession = &session.Session{
				ID:        uuid.New(),
				UserID:    userID,
				ExpiredAt: time.Now().Add(sessionTTL),
			}

			if err := s.sessionRepo.Create(ctx, userSession); err != nil {
				return fmt.Errorf("session_service.create_or_update %w", err)
			}
		} else {
			userSession.ExpiredAt = time.Now().Add(sessionTTL)
			if err := s.sessionRepo.Update(ctx, userSession); err != nil {
				return fmt.Errorf("session_service.create_or_update %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("session_service.create_or_update %w", err)
	}

	return nil
}
