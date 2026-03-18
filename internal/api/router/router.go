// Package router
package router

import (
	"net/http"

	"rest-api-app/internal/api/handlers"
	"rest-api-app/internal/api/handlers/students"
	"rest-api-app/internal/api/handlers/teachers"
)

type Handlers struct {
	Students *students.API
	Teachers *teachers.API
}

func Router(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	// mux.Handle("GET /", h.RootHandler)

	mux.Handle("GET /teachers/", h.Teachers.GetTeachersDB())
	mux.Handle("POST /teachers/", h.Teachers.AddTeacher())
	mux.Handle("PATCH /teachers/", h.Teachers.PatchTeachers())
	mux.Handle("DELETE /teachers/", h.Teachers.DeleteTeachers())
	//
	mux.Handle("GET /teachers/{id}", h.Teachers.GetTeacherByID())
	mux.Handle("PUT /teachers/{id}", h.Teachers.UpdateTeacher())
	mux.Handle("PATCH /teachers/{id}", h.Teachers.PatchOneTeacher())
	mux.Handle("DELETE /teachers/{id}", h.Teachers.DeleteTeacher())

	mux.HandleFunc("/execs/", handlers.ExecHandler)

	mux.Handle("GET /students/", h.Students.Get())
	mux.Handle("POST /students/", h.Students.AddStudent())
	mux.Handle("PATCH /students/", h.Students.PatchStudents())
	mux.Handle("DELETE /students/", h.Students.DeleteStudents())
	//
	mux.Handle("GET /students/{id}", h.Students.GetByID())
	mux.Handle("PUT /students/{id}", h.Students.UpdateStudent())
	mux.Handle("PATCH /students/{id}", h.Students.PatchOneStudent())
	mux.Handle("DELETE /students/{id}", h.Students.DeleteStudent())
	mux.Handle("GET /teachers/{id}/students", h.Teachers.GetStudentsByTeacher())

	return mux
}

// Why mux - mux позволяет определять много мульти роутов мульти эндпоинт
// Позволяет разделять логику
// Если у нас только несколько handlers можно без mux только http

// IN POSTMAN MUST HAVE LAST "/"!
// For example https://localhost:3000/students/
