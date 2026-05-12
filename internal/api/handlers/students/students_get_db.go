package students

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"
)

type GetStudentsDB interface {
	GetStudents(ctx context.Context, filters StudentFilters, page, limit int) ([]models.Student, int64, error)
}

func GetStudentsHandler(get GetStudentsDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filters := parseStudentFilters(r)
		slog.Debug("student filters", "filters", filters)

		page, limit := getPaginationParams(r)

		students, totalCount, err := get.GetStudents(r.Context(), filters, page, limit)
		if err != nil {
			slog.Error("get students", "err", err)
			http.Error(w, "Students not found", http.StatusNotFound)
			return
		}

		// url?limit=50&page=1
		// database will leave/will not show calculated entries from the beginning,
		// page -1 * limit ((1-1) * 50 = 0*50 = 0)
		// page 2, 2-1 *50 = 50, next 50 entries

		response := struct {
			Status   string           `json:"status"`
			Count    int64            `json:"count"`
			Page     int              `json:"page"`
			PageSize int              `json:"page_size"`
			Data     []models.Student `json:"data"`
		}{
			Status:   "success",
			Count:    totalCount,
			Page:     page,
			PageSize: limit,
			Data:     students,
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Debug("failed to encode response", "err", err)
			w.WriteHeader(http.StatusInternalServerError) // Send 500 status
			return
		}
	}
}

func getPaginationParams(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		slog.Debug("parsing page query parameter", "err", err)
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		slog.Debug("parsing limit query parameter", "err", err)
		limit = 10
	}
	return page, limit
}
