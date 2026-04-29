package handlers

import (
	"net/http"
	"rest-api-app/views/pages"
)

// func Make(h HTTPHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := h(w, r); err != nil {
// 			slog.Error("HTTP handler error", "err", err, "path", r.URL.Path)
// 		}
// 	}
// }

// func RegisterHandler(w http.ResponseWriter, r *http.Request) error {
// 	return pages.Register().Render(r.Context(), w)
// }

func RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pages.Register().Render(r.Context(), w)
	}
}
