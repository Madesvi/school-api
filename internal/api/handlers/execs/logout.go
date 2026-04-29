package execs

import (
	"net/http"
	"time"
)

func LogoutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     "Bearer",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			Expires:  time.Unix(0, 0),
			SameSite: http.SameSiteStrictMode,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "logout successfully"}`))
	}
}
