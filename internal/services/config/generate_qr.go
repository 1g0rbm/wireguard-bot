package config

import (
	"context"
	"fmt"

	"github.com/skip2/go-qrcode"
)

const (
	qrSize = 256
)

// GenerateQR method to generate vpn qr-config for userId.
func (s *ServiceConfig) GenerateQR(ctx context.Context, userID int64) ([]byte, error) {
	conf, err := s.GenerateConf(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate_qr] %w", err)
	}

	qrCode, err := qrcode.Encode(string(conf), qrcode.Medium, qrSize)
	if err != nil {
		return nil, fmt.Errorf("[config_service.generate_qr] %w", err)
	}

	return qrCode, nil
}
