package router

import (
	"net/http"
)

func teacherRoutes(mux *http.ServeMux, h Handlers) {
	mux.Handle("GET /teachers/", h.Teachers.GetTeachersDB())
	mux.Handle("POST /teachers/", h.Teachers.AddTeacher())
	mux.Handle("PATCH /teachers/", h.Teachers.PatchTeachers())
	mux.Handle("DELETE /teachers/", h.Teachers.DeleteTeachers())
	//
	mux.Handle("GET /teachers/{id}", h.Teachers.GetTeacherByID())
	mux.Handle("PUT /teachers/{id}", h.Teachers.UpdateTeacher())
	mux.Handle("PATCH /teachers/{id}", h.Teachers.PatchOneTeacher())
	mux.Handle("DELETE /teachers/{id}", h.Teachers.DeleteTeacher())
}
