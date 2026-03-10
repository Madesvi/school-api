package students

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type PatchOneStudent interface {
	PathOneStudent(ctx context.Context, id int, updates map[string]any) models.Student
}

func PatchOneStudentHandler(patchOne PatchOneStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("error convert to int")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updates map[string]any
		err = json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			log.Error().Err(err).Msg("error decoding request from body")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		existingStudent := patchOne.PathOneStudent(r.Context(), id, updates)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		err = json.NewEncoder(w).Encode(existingStudent)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
