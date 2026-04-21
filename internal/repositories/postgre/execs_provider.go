package postgre

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"rest-api-app/internal/api/handlers/execs"
	"rest-api-app/internal/models"
	"rest-api-app/pkg/utils"

	_ "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Лучше разделить ответственность postgreSQL и cache
// И собрать всё в Service
// TODO:
// Remove all logs! Add return fmt.Errorf("what error: %w", result.Error)
// Repo do not write logs

type ExecProvider struct {
	db *gorm.DB
}

func NewExecProvider(db *gorm.DB) *ExecProvider {
	return &ExecProvider{db: db}
}

func (p *ExecProvider) GetExecs(ctx context.Context, filters execs.ExecFilters) ([]models.Exec, error) {
	var exec []models.Exec
	tx := p.db.WithContext(ctx).Model(&exec)
	if filters.FirstName != "" {
		tx = tx.Where("first_name = ?", filters.FirstName)
	}
	if filters.LastName != "" {
		tx = tx.Where("last_name = ?", filters.LastName)
	}
	if filters.Email != "" {
		tx = tx.Where("email = ?", filters.Email)
	}
	if filters.UserName != "" {
		tx = tx.Where("username = ?", filters.UserName)
	}

	for _, sort := range filters.SortBy {
		tx = tx.Order(sort.Field + " " + sort.Order)
	}

	result := tx.Find(&exec)
	if result.Error != nil {
		return nil, result.Error
	}

	return exec, nil
}

func (p *ExecProvider) AddExec(ctx context.Context, newExec []models.Exec) ([]models.Exec, error) {
	if err := p.db.Create(newExec).Error; err != nil {
		return nil, fmt.Errorf("db create execs: %w", err)
	}
	return newExec, nil
}

func (p *ExecProvider) PatchExecs(ctx context.Context, updates []map[string]any) error {
	for _, update := range updates {
		idStr, ok := update["id"].(string)
		if !ok {
			return nil // CHECK!
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error().Err(err).Msg("invalid exec ID in update")
			// http.Error(w, "Error convert ID into int", http.StatusBadRequest)
		}

		gormUpdate := make(map[string]any)
		for k, v := range update {
			if k == "id" {
				continue
			}
			gormUpdate[utils.ToSnakeCase(k)] = v
		}

		result := p.db.Model(&models.Exec{}).Where("id = ?", id).Updates(gormUpdate)

		if result.Error != nil {
			log.Error().Err(result.Error).Msg("databese error during update")
			return result.Error
		}

		if result.RowsAffected == 0 {
			var count int64
			p.db.Model(&models.Exec{}).Where("id = ?", id).Count(&count)
			if count == 0 {
				return err
			}
		}
	}
	return nil
}

func (p *ExecProvider) GetOneExec(ctx context.Context, id int) (models.Exec, error) {
	var exec models.Exec
	result := p.db.First(&exec, id)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error")
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.Exec{}, result.Error
		}
		return models.Exec{}, result.Error
	}
	return exec, nil
}

func (p *ExecProvider) PathOneExec(ctx context.Context, id int, updates map[string]any) models.Exec {
	var existingExec models.Exec

	p.db.First(&existingExec, id)
	execVal := reflect.ValueOf(&existingExec).Elem() // Without &
	teacherType := execVal.Type()

	for k, v := range updates {
		for i := 0; i < execVal.NumField(); i++ {
			field := teacherType.Field(i)
			field.Tag.Get("json")
			if field.Tag.Get("json") == k+",omitempty" {
				if execVal.Field(i).CanSet() {
					execVal.Field(i).Set(reflect.ValueOf(v).Convert(execVal.Field(i).Type()))
				}
			}
		}
	}

	tx := p.db.Model(&models.Exec{})
	tx = tx.Where("id = ?", id)
	log.Info().Msgf("Id is: %d", id)
	tx.Updates(models.Exec{
		FirstName: existingExec.FirstName,
		LastName:  existingExec.LastName,
		Email:     existingExec.Email,
		UserName:  existingExec.UserName,
	})
	return existingExec
}

func (p *ExecProvider) DeleteOneExec(cxt context.Context, id int) error {
	var deleteExec models.Exec
	deleteExec.ID = id

	// If need to recieve rows deleted user
	result := p.db.Clauses(clause.Returning{}).Where("id = ?", id).Delete(&deleteExec)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("database error")
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return result.Error
	}

	// USE RowsAffected to check response from DB
	if result.RowsAffected == 0 {
		// http.Error(w, "Exec not found", http.StatusNotFound)
		return ErrNotFound
	}
	return nil
}
