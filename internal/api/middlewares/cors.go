package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
)

// CROSS-ORIGIN RESOURCE SHARING

// api is hosted at www.myapi.com
// fronted server is at www.myfronted.com

// Мы не проверяем DOMAIN! Мы проверяем Origin Header!

// Allowed origins
var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://www.myfronted.com",
	"http://localhost:3000",
	"http://localhost:7331",
}

func Cors(next http.Handler) http.Handler {
	fmt.Println("Cors Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Cors Middleware being returned...")
		origin := r.Header.Get("Origin")
		slog.Info("Origin Header", "origin", origin)
		path := r.URL.Path
		url := r.URL
		fmt.Printf("Path from r: %s and url: %s\n", path, url)

		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}

		// w.Header().Set()
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, HX-Request, HX-Target, HX-Current-URL")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// === Pre fligth check ===
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
		fmt.Println("Cors Middleware ends...")
	})
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
