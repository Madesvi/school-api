// Package teachers
package teachers

import (
	"context"
	"encoding/json"
	"net/http"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type TeacherGetDB interface {
	GetTeachers(ctx context.Context, filters TeacherFilters) ([]models.Teacher, error)
}

func GetTeachersHandler(get TeacherGetDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseTeacherFilters(r)
		teachers, err := get.GetTeachers(r.Context(), filters)
		log.Info().Msgf("Context: %v\n", r.Context())
		if err != nil {
			log.Error().Err(err).Msg("error")
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
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
			log.Error().Err(err).Msg("database error")
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}
	}
}
