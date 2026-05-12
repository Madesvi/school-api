package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"
)

//go:generate mockery
type GetStudentByID interface {
	GetStudentByID(ctx context.Context, id int) (models.Student, error)
}

func GetOneStudentHandler(get GetStudentByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Debug("convert to int", "err", err)
			http.Error(w, "Invalid ID format", http.StatusBadRequest) // 400
			return
		}

		student, err := get.GetStudentByID(r.Context(), id)
		if err != nil {
			slog.Warn("get one student", "err", err)
			http.Error(w, "Student not found", http.StatusNotFound) // 404
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(student)
		if err != nil {
			slog.Warn("encode student", "err", err)
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
