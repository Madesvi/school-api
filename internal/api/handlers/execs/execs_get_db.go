package execs

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"rest-api-app/internal/models"
)

type GetExecsDB interface {
	GetExecs(ctx context.Context, filters ExecFilters) ([]models.Exec, error)
}

func GetExecsHandler(get GetExecsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseExecFilters(r)
		execs, err := get.GetExecs(r.Context(), filters)
		if err != nil {
			slog.Error("failed to get execs", "err", err)
			http.Error(w, "Exec not found", http.StatusNotFound)
			return
		}

		response := struct {
			Status string        `json:"status"`
			Count  int           `json:"count"`
			Data   []models.Exec `json:"data"`
		}{
			Status: "success",
			Count:  len(execs),
			Data:   execs,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error encoding response", "err", err)
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}
	}
}
