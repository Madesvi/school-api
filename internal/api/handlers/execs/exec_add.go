// Package execs
package execs

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"

	"rest-api-app/internal/models"
	"rest-api-app/internal/services/mailer"
	"rest-api-app/pkg/utils"
)

type AddExec interface {
	AddExec(ctx context.Context, newExec []models.Exec) ([]models.Exec, error)
}

func GetFieldNames(model any) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd) // Get JSON tag
	}
	return fields // [id first_name last_name ...]
}

func AddExecHandler(add AddExec) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newExecs []models.Exec
		var rawExecs []map[string]any

		// WE JUST READ THE BODY ONCE - BYTE SLICE AND USE body for Unmarshal for more times
		// When we use r.Body more than once after first read our body will be empty
		// validator - LEARN GIT REPO HOW WE CAN USE!

		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error reading request body", "err", err)
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &rawExecs)
		if err != nil {
			slog.Error("failed to unmarshal body", "err", err)
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		// [id first_name last_name email username - password_changed_at password_reset_token password_code_expires role inactive_status created_at updated_at]
		fields := GetFieldNames(models.Exec{})
		// log.Info().Msgf("FIELDS: %v", fields)

		allowedFields := make(map[string]struct{})
		for _, field := range fields {
			allowedFields[field] = struct{}{}
		}

		for _, exec := range rawExecs {
			for key := range exec {
				_, ok := allowedFields[key]
				if !ok {
					http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
					return
				}
			}
		}

		err = json.Unmarshal(body, &newExecs)
		if err != nil {
			slog.Error("failed to unmarshal body", "err", err)
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		for i := range newExecs {
			// if exec.FirstName == "" || exec.LastName == "" || exec.Email == "" || exec.Class == "" || exec.Subject == "" {
			// 	http.Error(w, "All field are required", http.StatusBadRequest)
			// 	return
			// }
			encodedHash, err := utils.HashPassword(newExecs[i].Password)
			if err != nil {
				slog.Info("validation failed", "reason", "missing password", "user_email", newExecs[i].Email)
				http.Error(w, "Password is required", http.StatusBadRequest)
				return
			}

			newExecs[i].Password = encodedHash

			// if exec.Email == "" {
			// 	http...
			// }

			val := reflect.Indirect(reflect.ValueOf(newExecs[i]))
			for j := 0; j < val.NumField(); j++ {
				field := val.Field(j)
				if field.Kind() == reflect.String && field.String() == "" {
					http.Error(w, "All field are required", http.StatusBadRequest)
					return
				}
			}
		}

		newExecs, err = add.AddExec(r.Context(), newExecs)
		if err != nil {
			slog.Error("failed to add execs", "err", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// TEST SEND MAIL
		rawPass := "testtest" // For only test porposes
		if len(newExecs) > 0 {
			go func(u models.Exec, p string) {
				err := mailer.SendWelcomeEmail(newExecs[0].Email, "testtest")
				if err != nil {
					slog.Error("failed to send email", "err", err)
				}
			}(newExecs[0], rawPass)
		}
		// TEST SEND MAIL

		// Init var lastID for add value from DB
		var lastID int

		// Use range for read last ID
		for _, newExec := range newExecs {
			lastID = newExec.ID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status      string        `json:"status"`
			Count       int           `json:"count"`
			Data        []models.Exec `json:"data"`
			LastAddedID int           `json:"last_id"`
		}{
			Status:      "success",
			Count:       len(newExecs),
			Data:        newExecs,
			LastAddedID: lastID,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			slog.Error("error encoding response", "err", err)
		}
	}
}
