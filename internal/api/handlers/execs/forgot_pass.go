package execs

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"rest-api-app/internal/services/mailer"
	"strconv"
	"time"
)

// func Make(h HTTPHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := h(w, r); err != nil {
// 			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
// 		}
// 	}
// }

// func RootHandler(w http.ResponseWriter, r *http.Request) error {
// 	return Render(w, r, pages.Index())
// }

// func Router(h Handlers) *http.ServeMux {
// mux := http.NewServeMux()mux := http.NewServeMux()
// mux.Handle("GET /", handlers.Make(handlers.RootHandler))
// return mux

type ForgotPassword interface {
	CheckUserByEmail(ctx context.Context, email string) (int, error)
	UpdateTokenForOneExec(ctx context.Context, id int, token, exp string) error
}

func ForgotPasswordHandler(get ForgotPassword) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email string
		}
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invaid request body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		id, err := get.CheckUserByEmail(r.Context(), req.Email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}
		duration, err := strconv.Atoi(os.Getenv("RESET_TOKEN_EXP_DURATION"))
		if err != nil {
			slog.Error("Failed to parse reset token exp duration", "err", err)
			return
		}

		mins := time.Duration(duration)
		expiry := time.Now().Add(mins * time.Minute).Format(time.RFC3339)

		tokenBytes := make([]byte, 32)
		_, err = rand.Read(tokenBytes)
		if err != nil {
			slog.Error("Failed to generate reset token", "err", err)
			return
		}

		slog.Info("tokenBytes:", "token", tokenBytes)
		token := hex.EncodeToString(tokenBytes)
		slog.Info("token:", "token", token)

		hashedToken := sha256.Sum256(tokenBytes)
		slog.Info("hashedToken:", "token", hashedToken)

		hashedTokenString := hex.EncodeToString(hashedToken[:])

		err = get.UpdateTokenForOneExec(r.Context(), id, hashedTokenString, expiry)
		if err != nil {
			slog.Error("Failed to update token for one exec", "err", err)
			return
		}

		resetURL := fmt.Sprintf("http://localhost:3000/execs/resetpassword/reset/%s", token)
		message := fmt.Sprintf("Forgot your password? Reset your password using the following link:<br><a href='%s'>%s</a><br><br>If you didn't request a password reset, please ignore it. Link valid %d mins<br>", resetURL, resetURL, int(mins))

		go func() {
			err := mailer.SendWelcomeEmail(req.Email, message)
			if err != nil {
				slog.Error("failed to send email", "err", err)
			}
		}()

		fmt.Fprintf(w, "Password reset link sent to %s", req.Email)

	}
}
