package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserDoesNotHaveSession = errors.New("user does not have session")
	ErrExpiredSession         = errors.New("session expired")
)

func (s *ServiceSession) Check(ctx context.Context, sessionID uuid.UUID) error {
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("session_service.check.find_by_id %w", err)
	}

	if session == nil {
		return ErrUserDoesNotHaveSession
	}

	if session.ExpiredAt.Before(time.Now()) {
		return ErrExpiredSession
	}

	return nil
}

func (s *ServiceSession) CheckByUsername(ctx context.Context, username string) (*uuid.UUID, error) {
	session, err := s.sessionRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("session_service.check.find_by_username %w", err)
	}

	if session == nil {
		return nil, ErrUserDoesNotHaveSession
	}

	if session.ExpiredAt.Before(time.Now()) {
		return nil, ErrExpiredSession
	}

	return &session.ID, nil
}
