package students

import (
	"net/http"
	"strings"
)

// type Students struct {
// 	ID        int    `gorm:"primaryKey" json:"id"`
// 	FirstName string `gorm:"not null" json:"first_name,omitempty"`
// 	LastName  string `gorm:"not null" json:"last_name,omitempty"`
// 	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
// 	TeacherID int    `gorm:"not null" json:"teacher_id,omitempty"`
//
// 	Teacher *Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
//
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

type StudentFilters struct {
	FirstName string
	LastName  string
	Email     string
	SortBy    []SortField
}

type SortField struct {
	Field string
	Order string
}

func parseStudentFilters(r *http.Request) StudentFilters {
	var filters StudentFilters
	q := r.URL.Query()
	// log.Info().Msgf("Query: %v", q)

	if v := q.Get("first_name"); v != "" {
		filters.FirstName = v
	}
	if v := q.Get("last_name"); v != "" {
		filters.LastName = v
	}
	if v := q.Get("email"); v != "" {
		filters.Email = v
	}

	for _, param := range q["sortby"] {
		field, order, found := strings.Cut(param, ":")
		if !found {
			continue
		}
		// log.Info().Msgf("Field %v", field)
		// log.Info().Msgf("Order %v", order)
		// fmt.Println(isValidField(field))
		// fmt.Println(isValidOrder(order))
		if isValidField(field) && isValidOrder(order) {
			// log.Info().Msgf("Valid: %v, %v", field, order)
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
		"id":         true,
	}
	return allowed[field]
}

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}
