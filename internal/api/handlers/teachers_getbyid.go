package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type TeacherGetByID interface {
	GetTeacherByID(ctx context.Context, id int) (models.Teacher, error)
}

func GetOneTeacherHandler(get TeacherGetByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		teacher, err := get.GetTeacherByID(r.Context(), id)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(teacher)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
