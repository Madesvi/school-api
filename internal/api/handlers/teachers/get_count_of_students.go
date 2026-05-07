package teachers

import (
	"context"
	"encoding/json"
	"net/http"
	"rest-api-app/pkg/utils"
	"strconv"

	"github.com/rs/zerolog/log"
)

type GetCountStudent interface {
	GetStudentCountByTeacherID(ctx context.Context, teacherID int) (int64, error)
}

func GetStudentsCountHandler(get GetCountStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// admin, manager, exec
		_, err := utils.AuthorizeUser(r.Context().Value(utils.ContextKey("role")).(string), "admin", "manager", "exec")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		teacherIDStr := r.PathValue("id")
		teacherID, err := strconv.Atoi(teacherIDStr)
		if err != nil {
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		}

		// method from teacher_provider
		studentCount, err := get.GetStudentCountByTeacherID(r.Context(), teacherID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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
