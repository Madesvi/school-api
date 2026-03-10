// Package teachers
package teachers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
)

type DeleteTeachers interface {
	DeleteTeachers(ctx context.Context, ids []int) ([]int, error)
}

func DeleteTeachersHandler(delete DeleteTeachers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ids []int
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			log.Error().Err(err).Msg("cannot decode body")
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check provided IDs
		if len(ids) == 0 {
			http.Error(w, "No teacher IDs provided", http.StatusBadRequest)
		}

		deletedIDs, err := delete.DeleteTeachers(r.Context(), ids)
		if err != nil {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Status     string `json:"status"`
			DeletedIDs []int  `json:"deleted_ids"`
		}{
			Status:     "Teacher successfully deleted",
			DeletedIDs: deletedIDs,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
