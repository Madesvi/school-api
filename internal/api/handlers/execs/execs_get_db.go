package execs

import (
	"context"
	"encoding/json"
	"net/http"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type GetExecsDB interface {
	GetExecs(ctx context.Context, filters ExecFilters) ([]models.Exec, error)
}

func GetExecsHandler(get GetExecsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseExecFilters(r)
		log.Info().Msgf("Filters: %v", filters)
		execs, err := get.GetExecs(r.Context(), filters)
		if err != nil {
			log.Error().Err(err).Msg("error")
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
			log.Error().Err(err).Msg("database error")
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}
	}
}
