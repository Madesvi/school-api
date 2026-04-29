package execs

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"rest-api-app/internal/models"
	"rest-api-app/pkg/utils"
	"strconv"
	"time"
)

type UpdatePassword interface {
	CheckUser(ctx context.Context, id int) (string, string, string, error)
	UpdatePassForOneExec(ctx context.Context, id int, password string) error
}

func UpdatePasswordHandler(update UpdatePassword) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var req models.UpdatePasswordRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			slog.Info("Cannot decode", "err", err)
			http.Error(w, "Invalid Request Body", http.StatusBadRequest)
			return
		}
		r.Body.Close()

		if req.CurrentPassword == "" || req.NewPassword == "" {
			http.Error(w, "Please enter the password", http.StatusBadRequest)
			return
		}

		//	Search user if exists
		userName, userPassword, userRole, err := update.CheckUser(r.Context(), userId)
		if err != nil {
			slog.Warn("user check failed", "userId", userId, "err", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		err = utils.EncodeHash(req.CurrentPassword, userPassword)
		if err != nil {
			slog.Warn("password mismatch", "err", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			slog.Info("hashing failed", "err", err)
			http.Error(w, "Password is required", http.StatusInternalServerError)
			return
		}

		err = update.UpdatePassForOneExec(r.Context(), userId, hashedPassword)
		if err != nil {
			slog.Error("could not update the password", "err", err)
			http.Error(w, "Invalid credentials", http.StatusInternalServerError)
			return
		}

		token, err := utils.SignToken(idStr, userName, userRole)
		if err != nil {
			slog.Error("Failed to get token", "err", err)
			http.Error(w, "Could not create login token", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "Bearer",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Message string `json:"message"`
		}{
			Message: "Password updated",
		}
		json.NewEncoder(w).Encode(response)
	}
}
