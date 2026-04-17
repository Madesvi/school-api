// Package router
package router

import (
	"net/http"

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
