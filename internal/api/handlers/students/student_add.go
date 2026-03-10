// Package students
package students

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

type AddStudent interface {
	AddStudent(ctx context.Context, newStudent []models.Student) ([]models.Student, error)
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

// type Student struct {
// 	ID        int    `gorm:"primaryKey" json:"id"`
// 	FirstName string `gorm:"not null" json:"first_name,omitempty"`
// 	LastName  string `gorm:"not null" json:"last_name,omitempty"`
// 	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
// 	TeacherID int    `gorm:"not null" json:"teacher_id,omitempty"`
//
// 	Teacher *Teacher `gorm:"foreignKey:TeacherID" json:"teacher"`
//
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

func AddStudentHandler(add AddStudent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newStudents []models.Student
		var rawTeachers []map[string]any

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
			log.Info().Msg("Error here")
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		fields := GetFieldNames(models.Student{})

		allowedFields := make(map[string]struct{})
		for _, field := range fields {
			allowedFields[field] = struct{}{}
		}

		for _, student := range rawTeachers {
			for key := range student {
				_, ok := allowedFields[key]
				if !ok {
					http.Error(w, "Unacceptable field found in request. Only use allowed fields", http.StatusBadRequest)
					return
				}
			}
		}

		err = json.Unmarshal(body, &newStudents)
		if err != nil {
			http.Error(w, "Invalid request Body", http.StatusBadRequest)
			return
		}

		for _, student := range newStudents {
			// if student.FirstName == "" || student.LastName == "" || student.Email == "" || student.Class == "" || student.Subject == "" {
			// 	http.Error(w, "All field are required", http.StatusBadRequest)
			// 	return
			// }

			val := reflect.Indirect(reflect.ValueOf(student))
			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				if field.Kind() == reflect.String && field.String() == "" {
					http.Error(w, "All field are required", http.StatusBadRequest)
					return
				}
			}
		}

		newStudents, err = add.AddStudent(r.Context(), newStudents)
		if err != nil {
			log.Error().Err(err).Msg("error")
		}

		// Init var lastID for add value from DB
		var lastID int

		// Use range for read last ID
		for _, newStudent := range newStudents {
			lastID = newStudent.ID
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := struct {
			Status      string           `json:"status"`
			Count       int              `json:"count"`
			Data        []models.Student `json:"data"`
			LastAddedID int              `json:"last_id"`
		}{
			Status:      "success",
			Count:       len(newStudents),
			Data:        newStudents,
			LastAddedID: lastID,
		}
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error().Err(err).Msg("error encoding response")
		}
	}
}
