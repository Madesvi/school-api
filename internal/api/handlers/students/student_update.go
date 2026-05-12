package students

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"
)

type UpdateStudent interface {
	UpdateStudent(ctx context.Context, id int, updateStudent models.Student) (models.Student, error)
}

func UpdateStudentHandler(update UpdateStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Debug("cannot convert to int", "err", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateStudent models.Student

		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Debug("reading request body", "err", err)
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &updateStudent)
		if err != nil {
			slog.Debug("unnarshaling json", "err", err)
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}
		// err = json.NewDecoder(r.Body).Decode(&updateStudent)
		// if err != nil {
		// 	log.Error().Msg("error: ")
		// 	http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		// 	return
		// }

		updateStudentFromDB, err := update.UpdateStudent(r.Context(), id, updateStudent)
		if err != nil {
			slog.Warn("updating student", "err", err)
			http.Error(w, "Error updating student", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updateStudentFromDB)
		if err != nil {
			slog.Debug("encoding response", "err", err)
		}
	}
}
