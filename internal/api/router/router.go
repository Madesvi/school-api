// Package router
package router

import (
	"net/http"

	"rest-api-app/internal/api/handlers"
	"rest-api-app/internal/api/handlers/execs"
	"rest-api-app/internal/api/handlers/students"
	"rest-api-app/internal/api/handlers/teachers"
)

type Handlers struct {
	Teachers *teachers.API
	Students *students.API
	Execs    *execs.API
}

func Router(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./public"))
	mux.Handle("GET /public/", http.StripPrefix("/public/", fs))

	mux.Handle("GET /", handlers.Make(handlers.RootHandler))
	mux.Handle("GET /login", handlers.LoginHandler())
	mux.Handle("POST /login", handlers.LoginHandler())
	mux.Handle("GET /register", handlers.RegisterHandler())

	studentsRoutes(mux, h)
	teacherRoutes(mux, h)
	execsRoutes(mux, h)

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/students/
