// Package router
package router

import (
	"net/http"
)

type Handlers struct {
	RootHandler       http.Handler
	GetTeachers       http.Handler
	GetOneTeacher     http.Handler
	AddTeacher        http.Handler
	UpdateTeacher     http.Handler
	PatchTeachers     http.Handler
	PatchOneTeacher   http.Handler
	DeleteTeacherByID http.Handler
	DeleteTeachers    http.Handler

	GetStudents       http.Handler
	GetOneStudent     http.Handler
	AddStudent        http.Handler
	UpdateStudent     http.Handler
	PatchStudents     http.Handler
	PatchOneStudent   http.Handler
	DeleteStudentByID http.Handler
	DeleteStudents    http.Handler
}

func Router(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /", h.RootHandler)

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
	// mux.HandleFunc("/execs/", handlers.ExecHandler)

	mux.Handle("GET /students/", h.GetStudents)
	mux.Handle("POST /students/", h.AddStudent)
	mux.Handle("PATCH /students/", h.PatchStudents)
	mux.Handle("DELETE /students/", h.DeleteStudents)
	//
	mux.Handle("GET /students/{id}", h.GetOneStudent)
	mux.Handle("PUT /students/{id}", h.UpdateStudent)
	mux.Handle("PATCH /students/{id}", h.PatchOneStudent)
	mux.Handle("DELETE /students/{id}", h.DeleteStudentByID)

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/students/
