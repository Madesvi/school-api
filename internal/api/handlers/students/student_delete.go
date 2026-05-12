package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

type DeleteStudent interface {
	DeleteOneStudent(cxt context.Context, id int) error
}

func DeleteOneStudentHandler(delete DeleteStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Debug("convert to int", "err", err)
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		err = delete.DeleteOneStudent(r.Context(), id)
		if err != nil {
			slog.Warn("delete student", "err", err)
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
			slog.Debug("error encoding response", "err", err)
		}
	}
}
