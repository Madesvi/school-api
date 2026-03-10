package students

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type UpdateStudent interface {
	UpdateStudent(ctx context.Context, id int, updateStudent models.Student) (models.Student, error)
}

func UpdateStudentHandler(update UpdateStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var updateStudent models.Student

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(body, &updateStudent)
		log.Info().Msgf("Update student: %v", updateStudent)
		if err != nil {
			log.Error().Msg("error: ")
			http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
			return
		}
		// err = json.NewDecoder(r.Body).Decode(&updateStudent)
		// if err != nil {
		// 	log.Error().Msg("error: ")
		// 	http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		// 	return
		// }

		updateStudentFromDB, err := update.UpdateStudent(r.Context(), id, updateStudent)
		if err != nil {
			http.Error(w, "Error updating student", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(updateStudentFromDB)
		if err != nil {
			log.Error().Msg("Error")
		}
	}
}
