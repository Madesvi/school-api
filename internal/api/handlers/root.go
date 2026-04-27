// Package handlers
package handlers

import (
	"net/http"
	"rest-api-app/views/pages"
)

// func RootHandler() http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		foo.Index().Render(r.Context(), w)
// 		w.Write([]byte("Welcome to API!"))
// 	})
// }

func RootHandler(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, pages.Index())
}
