// Package handlers
package handlers

type Env struct {
	TeacherGetDB TeacherGetDB
}

func NewEnv(get TeacherGetDB) *Env {
	return &Env{TeacherGetDB: get}
}
