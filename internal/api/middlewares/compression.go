package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

// Compression use for minimize loading time for aplication
// Important for large styleshits, images, javascript files
// minimize data transfer

func Compression(next http.Handler) http.Handler {
	fmt.Println("Compression Middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check of the client accepts gzip encoding
		fmt.Println("Compression Middleware being returned...")
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
		} else {
			// Set the response Header
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()

			// Мы будем компресить Response который будет отправлен пользователю
			w = &gzipResponseWriter{ResponseWriter: w, Writer: gz}

			next.ServeHTTP(w, r)
			fmt.Println("Compression Middleware ends...")
		}
	})
}

// gzipResponseWriter wrap http.ResponseWriter to write gzipped response

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}
