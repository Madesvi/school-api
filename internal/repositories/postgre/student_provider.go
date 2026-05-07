package postgre

import (
	"context"
	"errors"
	"reflect"
	"strconv"

	"rest-api-app/internal/api/handlers/students"
	"rest-api-app/internal/models"
	"rest-api-app/pkg/utils"

	_ "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Лучше разделить ответственность postgreSQL и cache
// И собрать всё в Service

type StudentProvider struct {
	db *gorm.DB
}

func NewStudentProvider(db *gorm.DB) *StudentProvider {
	return &StudentProvider{db: db}
}

func (p *StudentProvider) GetStudents(ctx context.Context, filters students.StudentFilters, page, limit int) ([]models.Student, int64, error) {
	var student []models.Student
	var totalCount int64

	tx := p.db.WithContext(ctx).Model(&student)
	if filters.FirstName != "" {
		tx = tx.Where("first_name = ?", filters.FirstName)
	}
	if filters.LastName != "" {
		tx = tx.Where("last_name = ?", filters.LastName)
	}
	if filters.Email != "" {
		tx = tx.Where("email = ?", filters.Email)
	}

	if err := tx.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	for _, sort := range filters.SortBy {
		tx = tx.Order(sort.Field + " " + sort.Order)
	}

	offset := (page - 1) * limit

	result := tx.Limit(limit).Offset(offset).Find(&student)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return student, totalCount, nil
}

func (p *StudentProvider) GetStudentByID(ctx context.Context, id int) (models.Student, error) {
	var student models.Student
	result := p.db.First(&student, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Student{}, result.Error
		}
		return models.Student{}, result.Error
	}
	return student, nil
}

func (p *StudentProvider) AddStudent(ctx context.Context, newStudent []models.Student) ([]models.Student, error) {
	result := p.db.Create(newStudent)
	if result.Error != nil {
		log.Error().Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}

	return newStudent, nil
}

func (p *StudentProvider) UpdateStudent(ctx context.Context, id int, updateStudent models.Student) (models.Student, error) {
	updateStudent.ID = id

	if err := p.db.WithContext(ctx).Save(&updateStudent).Error; err != nil {
		log.Error().Err(err).Msg("err")
		return models.Student{}, err
	}

	return updateStudent, nil
}

func (p *StudentProvider) PatchStudents(ctx context.Context, updates []map[string]any) error {
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			return nil // CHECK!
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("invalid student ID in update")
			// http.Error(w, "Error convert ID into int", http.StatusBadRequest)
		}

		gormUpdate := make(map[string]any)
		for k, v := range update {
			if k == "id" {
				continue
			}
			gormUpdate[utils.ToSnakeCase(k)] = v
		}

		result := p.db.Model(&models.Student{}).Where("id = ?", id).Updates(gormUpdate)

		if result.Error != nil {
			log.Error().Err(result.Error).Msg("databese error during update")
			return result.Error
		}

		if result.RowsAffected == 0 {
			var count int64
			p.db.Model(&models.Student{}).Where("id = ?", id).Count(&count)
			if count == 0 {
				return err
			}
		}
	}
	return nil
}

func (p *StudentProvider) PathOneStudent(ctx context.Context, id int, updates map[string]any) models.Student {
	var existingStudent models.Student

	p.db.First(&existingStudent, id)
	studentVal := reflect.ValueOf(&existingStudent).Elem() // Without &
	teacherType := studentVal.Type()

	for k, v := range updates {
		for i := 0; i < studentVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if studentVal.Field(i).CanSet() {
					studentVal.Field(i).Set(reflect.ValueOf(v).Convert(studentVal.Field(i).Type()))
				}
			}
		}
	}

	tx := p.db.Model(&models.Student{})
	tx = tx.Where("id = ?", id)
	log.Info().Msgf("Id is: %d", id)
	tx.Updates(models.Student{
		FirstName: existingStudent.FirstName,
		LastName:  existingStudent.LastName,
		Email:     existingStudent.Email,
	})
	return existingStudent
}

func (p *StudentProvider) DeleteOneStudent(cxt context.Context, id int) error {
	var deleteStudent models.Student
	deleteStudent.ID = id

	// If need to recieve rows deleted user
	result := p.db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&deleteStudent)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return result.Error
	}

	// USE RowsAffected to check response from DB
	if result.RowsAffected == 0 {
		// http.Error(w, "Student not found", http.StatusNotFound)
		return ErrNotFound
	}
	return nil
}

func (p *StudentProvider) DeleteStudents(ctx context.Context, ids []int) ([]int, error) {
	var deleteStudent []models.Student
	// If need to recieve rows deleted user
	result := p.db.Clauses(clause.Returning{}).Delete(&deleteStudent, ids)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return nil, result.Error
	}

	deletedIDs := make([]int, len(deleteStudent)) // This is for len == to memory safe
	for i, student := range deleteStudent {
		deletedIDs[i] = student.ID
	}

	if len(deletedIDs) == 0 {
		return nil, ErrNotFound
	}
	return deletedIDs, nil
}
