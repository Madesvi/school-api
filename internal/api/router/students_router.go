package router

import "net/http"

func studentsRoutes(mux *http.ServeMux, h Handlers) {
	mux.Handle("GET /students/", h.Students.Get())
	mux.Handle("POST /students/", h.Students.AddStudent())
	mux.Handle("PATCH /students/", h.Students.PatchStudents())
	mux.Handle("DELETE /students/", h.Students.DeleteStudents())
	//
	mux.Handle("GET /students/{id}", h.Students.GetByID())
	mux.Handle("PUT /students/{id}", h.Students.UpdateStudent())
	mux.Handle("PATCH /students/{id}", h.Students.PatchOneStudent())
	mux.Handle("PATCH /students/payment/{id}", h.Students.Payment())
	mux.Handle("DELETE /students/{id}", h.Students.DeleteStudent())
	mux.Handle("GET /teachers/{id}/students", h.Teachers.GetStudentsByTeacher())
	mux.Handle("GET /teachers/{id}/studentscount", h.Teachers.GetCountStudent())
}
