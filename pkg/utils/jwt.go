package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userID, userName, role string) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtExpiresIn := os.Getenv("JWT_EXPIRES_IN")

	claims := jwt.MapClaims{
		"uid":  userID,
		"user": userName,
		"role": role,
	}

	if jwtExpiresIn == "" {
		duration, err := time.ParseDuration(jwtExpiresIn)
		if err != nil {
			return "", fmt.Errorf("cannot set expire: %w", err)
		}
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(duration))
	} else {
		claims["exp"] = jwt.NewNumericDate(time.Now().Add(15 * time.Minute))
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", fmt.Errorf("Internal error: %w", err)
	}

	return signedToken, nil
}
