package session

import (
	"wireguard-bot/internal/db"
	"wireguard-bot/internal/repository"
)

type ServiceSession struct {
	sessionRepo repository.SessionRepository
	txManager   db.TxManager
}

func NewServiceSession(
	sessionRepo repository.SessionRepository,
	txManager db.TxManager,
) *ServiceSession {
	return &ServiceSession{
		sessionRepo: sessionRepo,
		txManager:   txManager,
	}
}
