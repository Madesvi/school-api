package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type PathTeachers interface {
	PatchTeachers(ctx context.Context, updates []map[string]any) error
}

func PathTeachersHandler(patch PathTeachers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updates []map[string]any
		err := json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			log.Error().Err(err).Msg("error decoding response")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
		}
		log.Info().Msgf("Updates: %v", updates)

		err = patch.PatchTeachers(r.Context(), updates)
		if err != nil {
			log.Error().Err(err).Msg("error decoding response")
			http.Error(w, "Cannot update", http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
