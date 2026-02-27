// Package handlers provides handlers
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

func GetTeachersHandlers(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	if idStr == "" {

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
			log.Error().Msg("Error")
		}
		return
	}

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
		log.Error().Msg("Error")
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
		log.Error().Msg("Error")
	}
}

// PUT - complete replace
// PATH - modif some items

// PUT /teacher/{id}

func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error().Msg("error: ")
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// var updateTeacher models.Teacher
	// err = json.NewDecoder(r.Body).Decode(&updateTeacher)
	// if err != nil {
	// 	log.Error().Msg("error: ")
	// 	http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
	// 	return
	// }
	//
	// firstName := updateTeacher.FirstName
	// lastName := updateTeacher.LastName
	// email := updateTeacher.Email
	// class := updateTeacher.Class
	// subject := updateTeacher.Subject
	// // Find USER by ID
	// // Create instance of Teacher for record from DB
	// tx := postgre.DB.Model(&models.Teacher{})
	// tx = tx.Where("id = ?", id)
	// log.Info().Msgf("Id is: %d", id)
	// tx.Updates(models.Teacher{
	// 	FirstName: firstName,
	// 	LastName:  lastName,
	// 	Email:     email,
	// 	Class:     class,
	// 	Subject:   subject,
	// })

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	var updates map[string]any
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		log.Error().Msg("error: ")
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
	// response := struct {
	// 	Status       string `json:"status"`
	// 	LastUpdateID int    `json:"last_id"`
	// }{
	// 	Status:       "success",
	// 	LastUpdateID: id,
	// }

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	err = json.NewEncoder(w).Encode(existingTeacher)
	if err != nil {
		log.Error().Msg("Error")
	}

	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====
	// ===== REFACTOR PATCH =====

	// err = json.NewEncoder(w).Encode(updateTeacher)
	// if err != nil {
	// 	log.Error().Msg("Error")
	// }

	// postgre.DB.First(&existingTeacher, id)
	//
	// postgre.DB.Save(&updateTeacher{ID: idStr, first_name: "PEKIK"})
}

func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
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

func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// comment
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
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

// teacher/{id}
// teacger/9
// teacher/?key=value&query=value2&sortby=email&sortorder=ASC
