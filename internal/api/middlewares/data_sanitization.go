package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func XSSMiddleware(next http.Handler) http.Handler {
	slog.Info("XSS Middleware activated")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("XSS Middleware Ran")

		// Sanitize the URL Path
		sanitizePath, err := clean(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Info("Original Path", "Path:", r.URL.Path)
		slog.Info("Sanitize Path", "San Path:", sanitizePath)

		// Sanitize Query Params
		// For TEMPL and HTMX we dont't sanitize QUERY Params
		// HTML tags like <  > will be removed from the URL
		params := r.URL.Query()
		sanitizeQuery := make(map[string][]string)
		for key, values := range params {
			sanitizedKey, err := clean(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			var sanitizedValues []string
			for _, val := range values {
				cleanValue, err := clean(val)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				sanitizedValues = append(sanitizedValues, cleanValue.(string))
			}
			sanitizeQuery[sanitizedKey.(string)] = sanitizedValues
			fmt.Printf("Original Query %s: %s\n", key, strings.Join(values, ", "))
			fmt.Printf("Sanitized Query %s: %s\n", sanitizedKey, strings.Join(sanitizedValues, ", "))
		}

		r.URL.Path = sanitizePath.(string)
		r.URL.RawQuery = url.Values(sanitizeQuery).Encode()
		fmt.Println("Update URL:", r.URL.String())

		// Sanitize Body
		if r.Header.Get("Content-Type") == "application/json" {
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err != nil {
					slog.Error("Could not read request body", "err", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				bodyString := strings.TrimSpace(string(bodyBytes))
				slog.Info("Original Body", "Body", bodyString)

				// Reset the request body
				r.Body = io.NopCloser(bytes.NewReader([]byte(bodyString)))

				if len(bodyString) > 0 {
					var inputData any
					err := json.NewDecoder(bytes.NewReader([]byte(bodyString))).Decode(&inputData)
					if err != nil {
						http.Error(w, "Invalid JSON body", http.StatusBadRequest)
						return
					}
					slog.Info("Original JSON data", "data:", inputData)

					// Sanitize the JSON body
					sanitizedData, err := clean(inputData)
					if err != nil {
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}
					slog.Info("Sanitized Data", "data:", sanitizedData)

					// Marshaling the sanitized data to the body
					sanitizedBody, err := json.Marshal(sanitizedData)
					if err != nil {
						slog.Error("Error serializing data", "err", err)
						http.Error(w, err.Error(), http.StatusBadRequest)
						return
					}

					r.Body = io.NopCloser(bytes.NewReader(sanitizedBody))
					slog.Info("Sanitized body", "body", string(sanitizedBody))

				} else {
					slog.Info("Request is empty")
				}

			} else {
				slog.Warn("No body provided", "No body provided", r.Body)
			}
		} else if r.Header.Get("Content-Type") != "" {
			slog.Warn("Received request with unsupported Content-Type:", "Content-Type:", r.Header.Get("Content-Type"))
			http.Error(w, "Unsuppted Content-Type. Please use application/json", http.StatusUnsupportedMediaType)
			return
		}

		next.ServeHTTP(w, r)
		slog.Info("XSS Middleware Finished")
	})
}

// Clean sanitizes input to prevent XSS attacks. It removes any HTML
func clean(data any) (any, error) {
	switch v := data.(type) {
	case map[string]any:
		for key, val := range v {
			v[key] = sanitizeValue(val)
		}
		return v, nil
	case []any:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v, nil
	case string:
		return sanitizeString(v), nil

	default:
		return nil, fmt.Errorf("unsupported type: %T", data)
	}
}

func sanitizeValue(data any) any {
	switch v := data.(type) {
	case string:
		return sanitizeString(v)
	case map[string]any:
		for key, val := range v {
			v[key] = sanitizeValue(val)
		}
		return v
	case []any:
		for i, val := range v {
			v[i] = sanitizeValue(val)
		}
		return v
	default:
		return v
	}
}

func sanitizeString(value string) string {
	return bluemonday.UGCPolicy().Sanitize(value)
}
