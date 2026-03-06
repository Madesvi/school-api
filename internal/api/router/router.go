// Package router
package router

import (
	"net/http"
)

type Handlers struct {
	GetTeachers       http.Handler
	GetOneTeacher     http.Handler
	AddTeacher        http.Handler
	UpdateTeacher     http.Handler
	PatchTeachers     http.Handler
	PatchOneTeacher   http.Handler
	DeleteTeacherByID http.Handler
	DeleteTeachers    http.Handler
}

func Router(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	// mux.HandleFunc("/", handlers.RootHandler)

	mux.Handle("GET /teachers/", h.GetTeachers)
	mux.Handle("POST /teachers/", h.AddTeacher)
	mux.Handle("PATCH /teachers/", h.PatchTeachers)
	mux.Handle("DELETE /teachers/", h.DeleteTeachers)
	//
	mux.Handle("GET /teachers/{id}", h.GetOneTeacher)
	mux.Handle("PUT /teachers/{id}", h.UpdateTeacher)
	mux.Handle("PATCH /teachers/{id}", h.PatchOneTeacher)
	mux.Handle("DELETE /teachers/{id}", h.DeleteTeacherByID)
	//
	// mux.HandleFunc("/students/", handlers.StudentHandler)
	// mux.HandleFunc("/execs/", handlers.ExecHandler)

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/teachers/
