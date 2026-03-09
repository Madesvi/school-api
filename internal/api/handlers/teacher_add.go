package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
)

type AddTeacher interface {
	AddTeacher(ctx context.Context, newTeachers []models.Teacher) ([]models.Teacher, error)
}

func GetFieldNames(model any) []string {
	val := reflect.TypeOf(model)
	fields := []string{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldToAdd := strings.TrimSuffix(field.Tag.Get("json"), ",omitempty")
		fields = append(fields, fieldToAdd) // Get JSON tag
	}
	return fields
}

func AddTeacherHandler(add AddTeacher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newTeachers []models.Teacher
		var rawTeachers []map[string]interface{}

		// WE JUST READ THE BODY ONCE - BYTE SLICE AND USE body for Unmarshal for more times
		// When we use r.Body more than once after first read our body will be empty
		// validator - LEARN GIT REPO HOW WE CAN USE!

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		err = json.Unmarshal(body, &rawTeachers)
		if err != nil {
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		fields := GetFieldNames(models.Teacher{})

		allowedFields := make(map[string]struct{})
		for _, field := range fields {
			allowedFields[field] = struct{}{}
		}

		for _, teacher := range rawTeachers {
			for key := range teacher {
				_, ok := allowedFields[key]
				if !ok {
					http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
					return
				}
			}
		}

		err = json.Unmarshal(body, &newTeachers)
		if err != nil {
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		for _, teacher := range newTeachers {
			// if teacher.FirstName == "" || teacher.LastName == "" || teacher.Email == "" || teacher.Class == "" || teacher.Subject == "" {
			// 	http.Error(w, "All field are required", http.StatusBadRequest)
			// 	return
			// }

			val := reflect.Indirect(reflect.ValueOf(teacher))
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				if field.Kind() == reflect.String && field.String() == "" {
					http.Error(w, "All field are required", http.StatusBadRequest)
					return
				}
			}
		}

		newTeachers, err = add.AddTeacher(r.Context(), newTeachers)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		// Init var lastID for add value from DB
		var lastID int

		// Use range for read last ID
		for _, newTeacher := range newTeachers {
			lastID = newTeacher.ID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status      string           `json:"status"`
			Count       int              `json:"count"`
			Data        []models.Teacher `json:"data"`
			LastAddedID int              `json:"last_id"`
		}{
			Status:      "success",
			Count:       len(newTeachers),
			Data:        newTeachers,
			LastAddedID: lastID,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
