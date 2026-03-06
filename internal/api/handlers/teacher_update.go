package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type UpdateTeacher interface {
	UpdateTeacher(ctx context.Context, id int, updateTeacher models.Teacher) (models.Teacher, error)
}

func UpdateTeacherHandler(update UpdateTeacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateTeacher models.Teacher
		err = json.NewDecoder(r.Body).Decode(&updateTeacher)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		updateTeacherFromDB, err := update.UpdateTeacher(r.Context(), id, updateTeacher)
		if err != nil {
			http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updateTeacherFromDB)
		if err != nil {
			log.Error().Msg("Error")
		}
	}
}
