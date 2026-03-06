package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type AddTeacher interface {
	AddTeacher(ctx context.Context, newTeachers []models.Teacher) ([]models.Teacher, error)
}

func AddTeacherHandler(add AddTeacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTeachers []models.Teacher
		err := json.NewDecoder(r.Body).Decode(&newTeachers)
		if err != nil {
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		newTeachers, err = add.AddTeacher(r.Context(), newTeachers)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		// Init var lastID for add value from DB
		var lastID int

		// Use range for read last ID
		for _, newTeacher := range newTeachers {
			lastID = newTeacher.ID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status      string           `json:"status"`
			Count       int              `json:"count"`
			Data        []models.Teacher `json:"data"`
			LastAddedID int              `json:"last_id"`
		}{
			Status:      "success",
			Count:       len(newTeachers),
			Data:        newTeachers,
			LastAddedID: lastID,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
