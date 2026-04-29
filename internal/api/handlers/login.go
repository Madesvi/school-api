// Package handlers
package handlers

import (
	"net/http"
	"rest-api-app/views/components"
	"rest-api-app/views/pages"
	"strings"
)

// templ Login(vals struct{Email string}, errs LoginError) {

// func LoginHandler(w http.ResponseWriter, r *http.Request) error {
// 	return Render(w, r, pages.Login())
// }

func LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			pages.Login(struct{ Email string }{}, components.LoginError{}).Render(r.Context(), w)
			return
		}

		// TrimSpace - часто копируют с пробелами
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		errs := components.LoginError{}
		if email == "" {
			errs.Email = "Поле обязательно"
		}
		if password == "" {
			errs.Password = "Введите пароль"
		}

		if errs.Email != "" || errs.Password != "" {
			w.WriteHeader(http.StatusOK)
			components.LoginForm(struct{ Email string }{Email: email}, errs).Render(r.Context(), w)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
