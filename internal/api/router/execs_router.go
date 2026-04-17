package router

import "net/http"

func execsRoutes(mux *http.ServeMux, h Handlers) {
	mux.Handle("GET /execs/", h.Execs.Get())
	mux.Handle("POST /execs/", h.Execs.AddStudent())
	mux.Handle("PATCH /execs/", h.Execs.PatchStudents())

	mux.Handle("GET /execs/{id}", h.Execs.GetByID())
	mux.Handle("PATCH /execs/{id}", h.Execs.PatchOneStudent())
	mux.Handle("DELETE /execs/{id}", h.Execs.DeleteStudent())
	mux.Handle("POST /execs/{id}/updatepassword", h.Execs.DeleteStudent())

	mux.Handle("POST /execs/login", h.Execs.AddStudent())
	mux.Handle("POST /execs/logout", h.Execs.AddStudent())
	mux.Handle("POST /execs/forgotpassword", h.Execs.AddStudent())
	mux.Handle("POST /execs/resetpassword/reset/{resetcode}", h.Execs.AddStudent())
}
