package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"
)

type Payment interface {
	PaymentStudent(ctx context.Context, studentID int, fee int64) error
}

func PaymentHandler(pay Payment) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		idStr, err := strconv.Atoi(id)
		if err != nil {
			slog.Debug("convert to int", "err", err)
			http.Error(w, "Invalid ID format", http.StatusBadRequest) // 400
			return
		}

		var req struct {
			Fee int64 `json:"fee"`
		}

		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error().Err(err).Msg("error decoding request from body")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		if req.Fee < 0 {
			http.Error(w, "Fee must be greater than zero", http.StatusBadRequest)
			return
		}

		err = pay.PaymentStudent(r.Context(), idStr, req.Fee)
		if err != nil {
			slog.Warn("payment faild", "student_id", idStr, "err", err)
			http.Error(w, "Payment failed: execution error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		response := struct {
			Status  string `json:"status"`
			Payment int64  `json:"payment"`
		}{
			Status:  "success",
			Payment: (req.Fee / 100),
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Debug("encoding response", "err", err)
		}
	}
}
