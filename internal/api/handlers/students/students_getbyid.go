package students

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type StudentGetByID interface {
	GetStudentByID(ctx context.Context, id int) (models.Student, error)
}

func GetOneStudentHandler(get StudentGetByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		student, err := get.GetStudentByID(r.Context(), id)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(student)
		if err != nil {
			log.Error().Err(err).Msg("database error")
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}

		// result, err := json.Marshal(student)
		// if err != nil {
		// 	log.Error().Err(err).Msg("database error")
		// 	w.WriteHeader(http.StatusInternalServerError) // Send 500 status
		// 	return
		// }
		// if _, err := w.Write(result); err != nil {
		// 	log.Error().Err(err).Msg("failed to write response")
		// }
	}
}
