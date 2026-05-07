package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"rest-api-app/pkg/utils"

	"github.com/golang-jwt/jwt/v5"
)

//	TODO
//	Добавить структуру для хранения JWT секрета
// Или в main для передачи секрета в миддлвар напрямую
// os.Getenv - syscall = замедление

func JWTMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		slog.Debug("JWT Middleware is running")
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			slog.Debug("Inside JWT Middleware")
			token, err := r.Cookie("Bearer")
			if err != nil {
				slog.Debug("Authorization header not found")
				http.Error(w, "Authorization header not found", http.StatusUnauthorized)
				return
			}

			// jwtSecret := os.Getenv("JWT_SECRET")

			parsedToken, err := jwt.Parse(token.Value, func(token *jwt.Token) (any, error) {
				return secret, nil
			}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					slog.Debug("Token has expired")
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				} else if errors.Is(err, jwt.ErrTokenMalformed) {
					http.Error(w, "Token Malformed", http.StatusUnauthorized)
					return
				}
				slog.Error("Failed to parse token", "err", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if parsedToken.Valid {
				// slog.Debug("Valid token", "token", token.Value)
			} else {
				// slog.Debug("Invalid token", "token", token.Value)
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
			}
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				slog.Error("Invalid Login Token")
				http.Error(w, "Invalid Token", http.StatusUnauthorized)
				return
			}

			// ContextKey - свой тип для стринги
			ctx := context.WithValue(r.Context(), utils.ContextKey("role"), claims["role"])
			ctx = context.WithValue(ctx, utils.ContextKey("expiresAt"), claims["exp"])
			ctx = context.WithValue(ctx, utils.ContextKey("username"), claims["user"])
			ctx = context.WithValue(ctx, utils.ContextKey("userId"), claims["uid"])

			next.ServeHTTP(w, r.WithContext(ctx))
			slog.Debug("Exiting JWT Middleware")
		})
	}
}
