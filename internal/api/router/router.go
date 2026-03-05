// Package router
package router

import (
	"net/http"

	"rest-api-app/internal/api/handlers"
)

func Router(env *handlers.Env) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("GET /teachers/", env.GetTeachersHandler)
	// mux.HandleFunc("POST /teachers/", env.AddTeacherHandler)
	// mux.HandleFunc("PATCH /teachers/", env.PathTeachersHandlerNoReflection)
	// mux.HandleFunc("DELETE /teachers/", env.DeleteTeachersHandler)
	//
	// mux.HandleFunc("GET /teachers/{id}", env.GetOneTeacherHandler)
	// mux.HandleFunc("PUT /teachers/{id}", env.UpdateTeacherHandler)
	// mux.HandleFunc("PATCH /teachers/{id}", env.PatchOneTeacherHandler)
	// mux.HandleFunc("DELETE /teachers/{id}", env.DeleteOneTeacherHandler)
	//
	mux.HandleFunc("/students/", handlers.StudentHandler)
	mux.HandleFunc("/execs/", handlers.ExecHandler)

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/teachers/
