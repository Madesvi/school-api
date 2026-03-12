package teachers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type GetStudentsByTeacher interface {
	GetStudentsByTeacherID(ctx context.Context, teacherID int) ([]models.Student, error)
}

func GetStudentsByTeacherHendler(get GetStudentsByTeacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherIDStr := r.PathValue("id")
		teacherID, err := strconv.Atoi(teacherIDStr)
		if err != nil {
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		}

		// method from teacher_provider
		students, err := get.GetStudentsByTeacherID(r.Context(), teacherID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get students by teacher")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Responce
		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Student `json:"data"`
		}{
			Status: "success",
			Count:  len(students),
			Data:   students,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
