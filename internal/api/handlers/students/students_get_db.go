package students

import (
	"context"
	"encoding/json"
	"net/http"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type GetStudentsDB interface {
	GetStudents(ctx context.Context, filters StudentFilters) ([]models.Student, error)
}

func GetStudentsHandler(get GetStudentsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseStudentFilters(r)
		log.Info().Msgf("Filters: %v", filters)
		students, err := get.GetStudents(r.Context(), filters)
		// log.Info().Msgf("Context: %v\n", r.Context())
		if err != nil {
			log.Error().Err(err).Msg("error")
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Student `json:"data"`
		}{
			Status: "success",
			Count:  len(students),
			Data:   students,
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
