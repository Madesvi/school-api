package execs

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

//go:generate mockery
type GetOneExec interface {
	GetOneExec(ctx context.Context, id int) (models.Exec, error)
}

func GetOneExecHandler(get GetOneExec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("error")
			http.Error(w, "Invalid ID format", http.StatusBadRequest) // 400
			return
		}

		exec, err := get.GetOneExec(r.Context(), id)
		if err != nil {
			log.Error().Err(err).Msg("error")
			http.Error(w, "exec not found", http.StatusNotFound) // 404
			return
		}

		w.Header().Set("Content-Type", "application/json")

		err = json.NewEncoder(w).Encode(exec)
		if err != nil {
			log.Error().Err(err).Msg("database error")
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}

		// result, err := json.Marshal(exec)
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
