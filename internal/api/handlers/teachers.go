// Package handlers provides handlers
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"rest-api-app/internal/models"
	"rest-api-app/internal/repositories/postgre"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// teachers = make(map[int]models.Teacher)
var params = map[string]string{
	"first_name": "first_name",
	"last_name":  "last_name",
	"email":      "email",
	"class":      "class",
	"subject":    "subject",
}

// mu       = &sync.Mutex{}
// nextID = 1
//

func isValidOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	// Create instance of Teacher for record from DB
	teacherList := make([]models.Teacher, 0)
	// Take all rows from DB
	tx := postgre.DB.Model(&teacherList)

	for param := range params {
		value := r.URL.Query().Get(param)
		// log.Info().Msgf("Value from Query: %s", value)
		if value != "" {
			tx = tx.Where(param+" = ?", value)
		}
	}

	// URL /teacher/?sortby=name:asc&sortby=class:desc
	// Create sort
	sortParams := r.URL.Query()["sortby"]
	log.Info().Msgf("SortParams from query: %v", sortParams)
	if len(sortParams) > 0 {
		for _, param := range sortParams {
			parts := strings.Split(param, ":")
			log.Info().Msgf("Parts from query: %v", parts)
			if len(parts) != 2 {
				continue
			}
			// part[0] = name
			// part[1] = asc
			field, order := parts[0], parts[1]
			if !isValidField(field) || !isValidOrder(order) {
				continue
			}
			query := field + " " + order
			log.Info().Msgf("Applying sort: %s", query)
			tx.Order(query)
		}
	}
	// Get teachet or all teachers from our MODEL
	tx.Find(&teacherList)

	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}
}

func GetOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	// Find teacher from DB by the ID
	var teacher models.Teacher
	result := postgre.DB.First(&teacher, idStr)
	if result.Error != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(teacher)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// mu.Lock()
	// defer mu.Unlock()

	// Create new instance of Teachers for isnert from request json
	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		http.Error(w, "Invalid request Body", http.StatusBadRequest)
		return
	}

	// Use gorm for add many poerson to db (use Create)
	// postrge.DB.Create - where DB is a global variable for use in any handlers
	result := postgre.DB.Create(newTeachers)
	if result.Error != nil {
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return
		}
		return
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

func PathTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Error().Err(err).Msg("error decoding response")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			log.Error().Err(err).Msg("invalid teacher ID in update")
			http.Error(w, "Invalid teacher ID in update", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("invalid teacher ID in update")
			http.Error(w, "Error convert ID into int", http.StatusBadRequest)
		}

		var teacherFromDB models.Teacher
		result := postgre.DB.First(&teacherFromDB, id)
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("Teacher not found in DB")
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}

		// Apply update using reflaction
		teacherVal := reflect.ValueOf(&teacherFromDB).Elem() // Without &
		teacherType := teacherVal.Type()

		for k, v := range update {
			if k == "id" {
				continue // skip update id field
			}
			for i := 0; i < teacherVal.NumField(); i++ {
				field := teacherType.Field(i)
				if field.Tag.Get("json") == k+",omitempty" {
					fieldVal := teacherVal.Field(i)
					if fieldVal.CanSet() {
						val := reflect.ValueOf(v)
						if val.Type().ConvertibleTo(fieldVal.Type()) {
							fieldVal.Set(val.Convert(fieldVal.Type()))
						} else {
							log.Info().Msgf("cannot convert %v to %v", val.Type(), fieldVal.Type())
							return
						}
					}
					break
				}
			}
		}
		tx := postgre.DB.Model(&models.Teacher{})
		tx = tx.Where("id = ?", id)
		log.Info().Msgf("Id is: %d", id)
		result = tx.Updates(models.Teacher{
			FirstName: teacherFromDB.FirstName,
			LastName:  teacherFromDB.LastName,
			Email:     teacherFromDB.Email,
			Class:     teacherFromDB.Class,
			Subject:   teacherFromDB.Subject,
		})
		// Обработка ошибки если не удалост обновить
		if result.Error != nil {
			log.Error().Err(result.Error).Msg("database error")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func PathTeachersHandlerNoReflection(w http.ResponseWriter, r *http.Request) {
	var updates []map[string]any
	err := json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Error().Err(err).Msg("error decoding response")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			log.Error().Err(err).Msg("invalid teacher ID in update")
			http.Error(w, "Invalid teacher ID in update", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("invalid teacher ID in update")
			http.Error(w, "Error convert ID into int", http.StatusBadRequest)
		}

		gormUpdate := make(map[string]any)
		for k, v := range update {
			if k == "id" {
				continue
			}
			gormUpdate[toSnakeCase(k)] = v
		}

		result := postgre.DB.Model(&models.Teacher{}).Where("id = ?", id).Updates(gormUpdate)

		if result.Error != nil {
			log.Error().Err(result.Error).Msg("databese error during update")
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if result.RowsAffected == 0 {
			var count int64
			postgre.DB.Model(&models.Teacher{}).Where("id = ?", id).Count(&count)
			if count == 0 {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return
			}
		}

	}
	w.WriteHeader(http.StatusNoContent)
}

func PatchOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Err(err).Msg("error convert to int")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Error().Err(err).Msg("error decoding request from body")
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	postgre.DB.First(&existingTeacher, id)
	teacherVal := reflect.ValueOf(&existingTeacher).Elem() // Without &
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	tx := postgre.DB.Model(&models.Teacher{})
	tx = tx.Where("id = ?", id)
	log.Info().Msgf("Id is: %d", id)
	tx.Updates(models.Teacher{
		FirstName: existingTeacher.FirstName,
		LastName:  existingTeacher.LastName,
		Email:     existingTeacher.Email,
		Class:     existingTeacher.Class,
		Subject:   existingTeacher.Subject,
	})

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	// fmt.Println("TeacherType field 0: ", teacherVal.Type().Field(0))
	// fmt.Println("TeacherType field 1: ", teacherVal.Type().Field(1))

	// type Teacher struct {
	// 	ID        int    `gorm:"primaryKey" json:"id,omitempty"`
	// 	FirstName string `gorm:"not null" json:"first_name,omitempty"`
	// 	LastName  string `gorm:"not null" json:"last_name,omitempty"`
	// 	Email     string `gorm:"uniqueIndex;not null" json:"email,omitempty"`
	// 	Class     string `gorm:"not null" json:"class,omitempty"`
	// 	Subject   string `gorm:"not null" json:"subject,omitempty"`
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	// err = json.NewEncoder(w).Encode(updateTeacher)
	// if err != nil {
	// 	log.Error().Msg("Error")
	// }
}

func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Msg("error: ")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updateTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updateTeacher)
	if err != nil {
		log.Error().Msg("error: ")
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	tx := postgre.DB
	updateTeacher.ID = id

	if err := tx.Save(&updateTeacher).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(updateTeacher)
	if err != nil {
		log.Error().Msg("Error")
	}
}

func DeleteOneTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// comment
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Msg("error: ")
		http.Error(w, "Inavalid Request Payload", http.StatusBadRequest)
	}

	var deleteTeacher models.Teacher
	deleteTeacher.ID = id

	// If need to recieve rows deleted user
	result := postgre.DB.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&deleteTeacher)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// USE RowsAffected to check response from DB
	if result.RowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	// w.WriteHeader(http.StatusNoContent)
	// var deleteTeacher models.Teacher
	// postgre.DB.Delete(&deleteTeacher, id)

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted",
		ID:     id,
	}
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}
}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {
	var ids []int64
	err := json.NewDecoder(r.Body).Decode(&ids)
	if err != nil {
		log.Error().Err(err).Msg("cannot decode body")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Check provided IDs
	if len(ids) == 0 {
		http.Error(w, "No teacher IDs provided", http.StatusBadRequest)
	}

	var deleteTeacher []models.Teacher
	// If need to recieve rows deleted user
	result := postgre.DB.Clauses(clause.Returning{}).Delete(&deleteTeacher, ids)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	deletedIDs := make([]int, len(deleteTeacher)) // This is for len == to memory safe
	for i, teacher := range deleteTeacher {
		deletedIDs[i] = teacher.ID
	}

	if len(deletedIDs) == 0 {
		http.Error(w, "No teacher found to delete", http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	response := struct {
		Status     string `json:"status"`
		DeletedIDs []int  `json:"deleted_ids"`
	}{
		Status:     "Teacher successfully deleted",
		DeletedIDs: deletedIDs,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error().Err(err).Msg("error encoding response")
	}
}
