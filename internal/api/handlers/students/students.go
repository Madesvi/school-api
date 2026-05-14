package students

import "net/http"

type StudentsProvider interface {
	AddStudent
	DeleteStudents
	PatchOneStudent
	UpdateStudent
	DeleteStudent
	GetStudentsDB
	GetStudentByID
	PatchStudents
	Payment
}

type API struct{ p StudentsProvider }

func NewAPI(p StudentsProvider) *API { return &API{p: p} }

func (a *API) AddStudent() http.HandlerFunc      { return AddStudentHandler(a.p) }
func (a *API) DeleteStudents() http.HandlerFunc  { return DeleteStudentsHandler(a.p) }
func (a *API) PatchOneStudent() http.HandlerFunc { return PatchOneStudentHandler(a.p) }
func (a *API) UpdateStudent() http.HandlerFunc   { return UpdateStudentHandler(a.p) }
func (a *API) DeleteStudent() http.HandlerFunc   { return DeleteOneStudentHandler(a.p) }
func (a *API) Get() http.HandlerFunc             { return GetStudentsHandler(a.p) }
func (a *API) GetByID() http.HandlerFunc         { return GetOneStudentHandler(a.p) }
func (a *API) PatchStudents() http.HandlerFunc   { return PatchStudentsHandler(a.p) }
func (a *API) Payment() http.HandlerFunc         { return PaymentHandler(a.p) }
