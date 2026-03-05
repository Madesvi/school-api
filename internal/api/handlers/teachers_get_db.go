package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type TeacherGetDB interface {
	GetTeachers(ctx context.Context, filters TeacherFilters) ([]models.Teacher, error)
}

type TeacherFilters struct {
	FirstName string
	LastName  string
	Email     string
	Class     string
	Subject   string
	SortBy    []SortField
}

type SortField struct {
	Field string
	Order string
}

func parseTeacherFilters(r *http.Request) TeacherFilters {
	var filters TeacherFilters
	q := r.URL.Query()

	if v := q.Get("first_name"); v != "" {
		filters.FirstName = v
	}
	if v := q.Get("last_name"); v != "" {
		filters.LastName = v
	}
	if v := q.Get("email"); v != "" {
		filters.Email = v
	}
	if v := q.Get("class"); v != "" {
		filters.Class = v
	}
	if v := q.Get("subject"); v != "" {
		filters.Subject = v
	}

	for _, param := range q["sortby"] {
		field, order, found := strings.Cut(param, ":")
		if !found {
			continue
		}
		if isValidField(field) && isValidOrder(order) {
			filters.SortBy = append(filters.SortBy, SortField{Field: field, Order: order})
		}
	}

	return filters
}

func isValidField(field string) bool {
	allowed := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
		"id":         true,
	}
	return allowed[field]
}

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func (e *Env) GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	filters := parseTeacherFilters(r)
	teachers, err := e.TeacherGetDB.GetTeachers(r.Context(), filters)
	if err != nil {
		log.Error().Err(err).Msg("error")
	}

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teachers,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}
}
