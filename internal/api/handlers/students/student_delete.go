package students

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type DeleteStudent interface {
	DeleteOneStudent(cxt context.Context, id int) error
}

func DeleteOneStudentHandler(delete DeleteStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		}

		err = delete.DeleteOneStudent(r.Context(), id)
		if err != nil {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		response := struct {
			Status string `json:"status"`
			ID     int    `json:"id"`
		}{
			Status: "Student successfully deleted",
			ID:     id,
		}
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
