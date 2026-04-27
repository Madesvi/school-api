// Package handlers
package handlers

import (
	"net/http"
	"rest-api-app/views/pages"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, pages.Login())
}
