package execs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"rest-api-app/internal/models"
	"rest-api-app/pkg/utils"
	"time"
)

type ResetPassword interface {
	GetResetToken(ctx context.Context, resetToken, exp string) (models.Exec, error)
	UpdatePassForOneExec(ctx context.Context, id int, password string) error
}

func ResetPasswordHandler(reset ResetPassword) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.PathValue("resetcode")

		type request struct {
			NewPassword     string `json:"new_password"`
			ConfirmPassword string `json:"confirm_password"`
		}

		var req request

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid values in request", http.StatusBadRequest)
		}

		if req.NewPassword != req.ConfirmPassword {
			http.Error(w, "Password should match", http.StatusBadRequest)
			return
		}

		bytes, err := hex.DecodeString(token)
		if err != nil {
			slog.Error("Failed to decode reset code", "err", err)
			return
		}

		hashedToken := sha256.Sum256(bytes)
		hashedTokenString := hex.EncodeToString(hashedToken[:])

		user, err := reset.GetResetToken(r.Context(), hashedTokenString, time.Now().Format(time.RFC3339))
		if err != nil {
			slog.Error("Invalid or expire reset code", "err", err)
			http.Error(w, "Invalid or expire reset code", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			slog.Error("Internal error", "err", err)
			return
		}

		err = reset.UpdatePassForOneExec(r.Context(), user.ID, hashedPassword)
		if err != nil {
			slog.Error("Failed to update token for one exec", "err", err)
			return
		}
		fmt.Fprintln(w, "Password reset successfully")
	}
}
