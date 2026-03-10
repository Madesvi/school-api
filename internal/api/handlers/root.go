// Package handlers
package handlers

import "net/http"

func RootHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintf(w, "Hello Root Route")
		w.Write([]byte("Welcome to API"))
	})
}
