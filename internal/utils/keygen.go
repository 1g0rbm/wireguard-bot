package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func GenerateKeyPair() (string, string, error) {
	var privateKey [32]byte
	_, err := rand.Read(privateKey[:])
	if err != nil {
		return "", "", fmt.Errorf("utils.generate_key_pair %w", err)
	}

	var publicKey [32]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)

	privateKeyStr := base64.StdEncoding.EncodeToString(privateKey[:])
	publicKeyStr := base64.StdEncoding.EncodeToString(publicKey[:])

	return privateKeyStr, publicKeyStr, nil
}
