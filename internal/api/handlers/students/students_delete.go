package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type DeleteStudents interface {
	DeleteStudents(ctx context.Context, ids []int) ([]int, error)
}

func DeleteStudentsHandler(delete DeleteStudents) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ids []int
		err := json.NewDecoder(r.Body).Decode(&ids)
		if err != nil {
			slog.Debug("decode body", "err", err)
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Check provided IDs
		if len(ids) == 0 {
			http.Error(w, "No student IDs provided", http.StatusBadRequest)
			return
		}

		deletedIDs, err := delete.DeleteStudents(r.Context(), ids)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Status     string `json:"status"`
			DeletedIDs []int  `json:"deleted_ids"`
		}{
			Status:     "Student successfully deleted",
			DeletedIDs: deletedIDs,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Warn("encoding response", "err", err)
		}
	}
}
