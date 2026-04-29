package router

import "net/http"

func execsRoutes(mux *http.ServeMux, h Handlers) {
	mux.Handle("GET /execs/", h.Execs.GetExecs())
	mux.Handle("POST /execs/", h.Execs.AddExec())
	mux.Handle("PATCH /execs/", h.Execs.PatchExecs())

	mux.Handle("GET /execs/{id}", h.Execs.GetOneExec())
	mux.Handle("PATCH /execs/{id}", h.Execs.PatchOneExec())
	mux.Handle("DELETE /execs/{id}", h.Execs.DeleteOneExec())
	mux.Handle("POST /execs/{id}/updatepassword", h.Execs.UpdatePassword())

	mux.Handle("POST /execs/login", h.Execs.LoginUser())
	mux.Handle("POST /execs/logout", h.Execs.LogoutUser())
	mux.Handle("POST /execs/forgotpassword", h.Execs.GetExecs())
	mux.Handle("POST /execs/resetpassword/reset/{resetcode}", h.Execs.GetExecs())
}
