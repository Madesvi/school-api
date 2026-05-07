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

type LoginUser interface {
	LoginUser(ctx context.Context, username string) (models.Exec, error)
}

func LoginUserHandler(login LoginUser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req models.Exec
		//	Data validation
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
		}
		defer r.Body.Close()

		if req.UserName == "" || req.Password == "" {
			http.Error(w, "Username and password are requrired", http.StatusBadRequest)
			return
		}

		//	Search user if exists
		user, err := login.LoginUser(r.Context(), req.UserName)
		if err != nil {
			slog.Error("user not found", "err", err)
			http.Error(w, "Invalid credentials", http.StatusBadRequest)
			return
		}
		//	Is user active
		if user.InactiveStatus {
			http.Error(w, "User is inactive", http.StatusForbidden)
			return
		}

		//	Verify password
		err = utils.EncodeHash(req.Password, user.Password)
		if err != nil {
			slog.Error("cannot encode password", "err", err)
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		//	Generate token
		tokenString, err := utils.SignToken(strconv.Itoa(user.ID), user.UserName, user.Role)
		if err != nil {
			slog.Error("Failed to get token", "err", err)
			http.Error(w, "Could not create login token", http.StatusInternalServerError)
			return
		}

		//	Send token as a response or as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "Bearer",
			Value:    tokenString,
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteStrictMode, // Prevents cross site attack
		})
		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Token string `json:"token"`
		}{
			Token: tokenString,
		}
		json.NewEncoder(w).Encode(response)
	}
}
