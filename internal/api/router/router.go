package router

import (
	"net/http"

	"rest-api-app/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/teachers/", handlers.TeacherHandler)
	mux.HandleFunc("/students/", handlers.StudentHandler)
	mux.HandleFunc("/execs/", handlers.ExecHandler)

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/teachers/
