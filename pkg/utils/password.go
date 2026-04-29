package utils

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"log/slog"
	"strings"

	"golang.org/x/crypto/argon2"
)

func EncodeHash(password, encodedHash string) error {
	//	Verify password
	parts := strings.Split(encodedHash, ".")
	if len(parts) != 2 {
		return errors.New("invalid encode hash format")
	}

	saltBase64 := parts[0]
	hashedPasswordBase64 := parts[1]

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		slog.Error("Failed to decode salt", "err", err)
		return err
	}
	hashedPassword, err := base64.StdEncoding.DecodeString(hashedPasswordBase64)
	if err != nil {
		slog.Error("Failed to decode hashed password", "err", err)
		return err
	}

	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	if len(hash) != len(hashedPassword) {
		slog.Error("hash length mismatch", "err", err)
		return err
	}

	if subtle.ConstantTimeCompare(hash, hashedPassword) == 1 {
		return nil
	}
	return errors.New("password does not match")
}
