package postgre

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"rest-api-app/internal/models"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

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

func addFilters(tx *gorm.DB, r *http.Request) *gorm.DB {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param := range params {
		value := r.URL.Query().Get(param)
		fmt.Printf("VALUE: %s PARAM: %s\n", value, param)
		if value != "" {
			tx = tx.Where(param+" = ?", value)
		}
	}
	return tx
}

func addSorting(tx *gorm.DB, r *http.Request) *gorm.DB {
	sortParams := r.URL.Query()["sortby"]
	log.Info().Msgf("SortParams from query: %v", sortParams)
	if len(sortParams) > 0 {
		for _, param := range sortParams {
			// parts := strings.Split(param, ":")
			// // log.Info().Msgf("Parts from query: %v", parts)
			// if len(parts) != 2 {
			// 	continue
			// }
			// part[0] = name
			// part[1] = asc
			// field, order := parts[0], parts[1]
			field, order, found := strings.Cut(param, ":")
			if !found {
				continue
			}
			if !isValidField(field) || !isValidOrder(order) {
				continue
			}
			query := field + " " + order
			log.Info().Msgf("Applying sort: %s", query)
			tx.Order(query)
		}
	}
	return tx
}

func GetTeachersDBHandler(teachers []models.Teacher, r *http.Request) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to db")
	}
	log.Info().Msg("Connect to PostgreSQL DB")

	tx := db.Model(&teachers)
	tx = addFilters(tx, r)
	tx = addSorting(tx, r)
	tx.Find(&teachers)

	return teachers, nil
}

func GetTeacherByID(w http.ResponseWriter, id int) (models.Teacher, bool) {
	// Find teacher from DB by the ID
	var teacher models.Teacher
	db, err := ConnectDB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to db")
		return teacher, true
	}
	log.Info().Msg("Connect to PostgreSQL DB")

	result := db.First(&teacher, id)
	if result.Error != nil {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return teacher, true
		}
		return teacher, true
	}
	return teacher, false
}
