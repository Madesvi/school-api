package postgre

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"rest-api-app/internal/api/handlers"
	"rest-api-app/internal/models"

	_ "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type teacherProvider struct {
	db *gorm.DB
	// redis *redis.Client
}

func NewTeacherProvider(db *gorm.DB) *teacherProvider {
	return &teacherProvider{db: db}
}

var ErrNotFound = errors.New("record not found")

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

func (p *teacherProvider) GetTeachers(ctx context.Context, filters handlers.TeacherFilters) ([]models.Teacher, error) {
	var teachers []models.Teacher
	tx := p.db.Model(&teachers)
	if filters.FirstName != "" {
		tx = tx.Where("first_name = ?", filters.FirstName)
	}
	if filters.LastName != "" {
		tx = tx.Where("last_name = ?", filters.LastName)
	}
	// ... ADD ALL FILTERS!!!

	for _, sort := range filters.SortBy {
		tx = tx.Order(sort.Field + " " + sort.Order)
	}

	err := tx.Find(&teachers).Error

	return teachers, err
}

func GetTeacherByID(db *gorm.DB, id int) (models.Teacher, error) {
	// Find teacher from DB by the ID
	var teacher models.Teacher

	result := db.First(&teacher, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Teacher{}, result.Error
		}
		return models.Teacher{}, result.Error
	}
	return teacher, nil
}

func AddTeacher(db *gorm.DB, newTeachers []models.Teacher) ([]models.Teacher, error) {
	result := db.Create(newTeachers)
	if result.Error != nil {
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return newTeachers, nil
}

func UpdateTeacher(db *gorm.DB, id int, updateTeacher models.Teacher) (models.Teacher, error) {
	// tx := postgre.DB
	updateTeacher.ID = id

	if err := db.Save(&updateTeacher).Error; err != nil {
		log.Error().Err(err).Msg("err")
		return models.Teacher{}, err
	}
	return updateTeacher, nil
}

func PatchTeachers(db *gorm.DB, updates []map[string]any) error {
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			return nil // CHECK!
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("invalid teacher ID in update")
			// http.Error(w, "Error convert ID into int", http.StatusBadRequest)
		}

		gormUpdate := make(map[string]any)
		for k, v := range update {
			if k == "id" {
				continue
			}
			gormUpdate[toSnakeCase(k)] = v
		}

		result := db.Model(&models.Teacher{}).Where("id = ?", id).Updates(gormUpdate)

		if result.Error != nil {
			log.Error().Err(result.Error).Msg("databese error during update")
			// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return result.Error
		}

		if result.RowsAffected == 0 {
			var count int64
			db.Model(&models.Teacher{}).Where("id = ?", id).Count(&count)
			if count == 0 {
				// http.Error(w, "Teacher not found", http.StatusNotFound)
				return err
			}
		}
	}
	return nil
}

func PathOneTeacher(db *gorm.DB, id int, updates map[string]any) models.Teacher {
	var existingTeacher models.Teacher

	db.First(&existingTeacher, id)
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

	tx := db.Model(&models.Teacher{})
	tx = tx.Where("id = ?", id)
	log.Info().Msgf("Id is: %d", id)
	tx.Updates(models.Teacher{
		FirstName: existingTeacher.FirstName,
		LastName:  existingTeacher.LastName,
		Email:     existingTeacher.Email,
		Class:     existingTeacher.Class,
		Subject:   existingTeacher.Subject,
	})
	return existingTeacher
}

func DeleteOneTeacher(db *gorm.DB, id int) error {
	var deleteTeacher models.Teacher
	deleteTeacher.ID = id

	// If need to recieve rows deleted user
	result := db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&deleteTeacher)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return result.Error
	}

	// USE RowsAffected to check response from DB
	if result.RowsAffected == 0 {
		// http.Error(w, "Teacher not found", http.StatusNotFound)
		return ErrNotFound
	}
	return nil
}

func DeleteTeachers(db *gorm.DB, ids []int) ([]int, error) {
	var deleteTeacher []models.Teacher
	// If need to recieve rows deleted user
	result := db.Clauses(clause.Returning{}).Delete(&deleteTeacher, ids)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, result.Error
	}

	deletedIDs := make([]int, len(deleteTeacher)) // This is for len == to memory safe
	for i, teacher := range deleteTeacher {
		deletedIDs[i] = teacher.ID
	}

	if len(deletedIDs) == 0 {
		return nil, ErrNotFound
	}
	return deletedIDs, nil
}
