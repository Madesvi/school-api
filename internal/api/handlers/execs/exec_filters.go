package execs

import (
	"net/http"
	"strings"
)

type ExecFilters struct {
	FirstName string
	LastName  string
	Email     string
	UserName  string
	SortBy    []SortField
}

type SortField struct {
	Field string
	Order string
}

func parseExecFilters(r *http.Request) ExecFilters {
	var filters ExecFilters
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
	if v := q.Get("username"); v != "" {
		filters.UserName = v
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
		"username":   true,
		"id":         true,
	}
	return allowed[field]
}

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}
