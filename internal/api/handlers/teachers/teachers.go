package teachers

import "net/http"

type TeachersProvider interface {
	AddTeacher
	DeleteTeachers
	PatchOneTeacher
	UpdateTeacher
	DeleteTeacher
	GetTeachersDB
	GetTeacherByID
	PatchTeachers
}

type API struct{ p TeachersProvider }

func NewAPI(p TeachersProvider) *API { return &API{p: p} }

func (a *API) AddTeacher() http.HandlerFunc      { return AddTeacherHandler(a.p) }
func (a *API) DeleteTeachers() http.HandlerFunc  { return DeleteTeachersHandler(a.p) }
func (a *API) PatchOneTeacher() http.HandlerFunc { return PatchOneTeacherHandler(a.p) }
func (a *API) UpdateTeacher() http.HandlerFunc   { return UpdateTeacherHandler(a.p) }
func (a *API) DeleteTeacher() http.HandlerFunc   { return DeleteOneTeacherHandler(a.p) }
func (a *API) GetTeachersDB() http.HandlerFunc   { return GetTeachersHandler(a.p) }
func (a *API) GetTeacherByID() http.HandlerFunc  { return GetOneTeacherHandler(a.p) }
func (a *API) PatchTeachers() http.HandlerFunc   { return PatchOneTeacherHandler(a.p) }
