package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"
)

type PatchOneStudent interface {
	PathOneStudent(ctx context.Context, id int, updates map[string]any) models.Student
}

func PatchOneStudentHandler(patchOne PatchOneStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Debug("convert to int", "err", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updates map[string]any
		err = json.NewDecoder(r.Body).Decode(&updates)
		if err != nil {
			slog.Debug("decoding request from body", "err", err)
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		existingStudent := patchOne.PathOneStudent(r.Context(), id, updates)
		currentBalance := existingStudent.Balance / 100

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		response := struct {
			Status  string         `json:"status"`
			Balance int64          `josn:"balance"`
			Data    models.Student `json:"data"`
		}{
			Status:  "success",
			Balance: currentBalance,
			Data:    existingStudent,
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Debug("encoding response", "err", err)
		}
	}
}
