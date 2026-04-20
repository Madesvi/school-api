package execs

import "net/http"

type ExecsProvider interface {
	GetExecsDB
	AddExec
	PatchExecs
	GetOneExec
	PatchOneExec
	DeleteOneExec
}

type API struct{ p ExecsProvider }

func NewAPI(p ExecsProvider) *API { return &API{p: p} }

func (a *API) GetExecs() http.HandlerFunc      { return GetExecsHandler(a.p) }
func (a *API) AddExec() http.HandlerFunc       { return AddExecHandler(a.p) }
func (a *API) PatchExecs() http.HandlerFunc    { return PatchExecsHandler(a.p) }
func (a *API) GetOneExec() http.HandlerFunc    { return GetOneExecHandler(a.p) }
func (a *API) PatchOneExec() http.HandlerFunc  { return PatchOneExecHandler(a.p) }
func (a *API) DeleteOneExec() http.HandlerFunc { return DeleteOneExecHandler(a.p) }
