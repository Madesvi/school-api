package teachers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type GetCountStudent interface {
	GetStudentCountByTeacherID(ctx context.Context, teacherID int) (int64, error)
}

func GetStudentsCountHandler(get GetCountStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		teacherIDStr := r.PathValue("id")
		teacherID, err := strconv.Atoi(teacherIDStr)
		if err != nil {
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		}

		// method from teacher_provider
		studentCount, err := get.GetStudentCountByTeacherID(r.Context(), teacherID)
		if err != nil {
			return
		}

		// Responce
		response := struct {
			Status string `json:"status"`
			Count  int64  `json:"count"`
		}{
			Status: "success",
			Count:  studentCount,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
