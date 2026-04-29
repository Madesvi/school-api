package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is blank")
	}
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		slog.Error("failed to generate salt", "err", err)
		return "", fmt.Errorf("failed to generate salt %w", err)
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashBase64 := base64.StdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf("%s.%s", saltBase64, hashBase64)
	return encodedHash, nil
}
