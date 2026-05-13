package postgre

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"

	"rest-api-app/internal/api/handlers/students"
	"rest-api-app/internal/models"
	"rest-api-app/pkg/utils"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Лучше разделить ответственность postgreSQL и cache
// И собрать всё в Service

type StudentProvider struct {
	db  *gorm.DB
	pgx *pgxpool.Pool
}

func NewStudentProvider(db *gorm.DB, pgx *pgxpool.Pool) *StudentProvider {
	return &StudentProvider{db: db, pgx: pgx}
}

func (p *StudentProvider) GetStudents(ctx context.Context, filters students.StudentFilters, page, limit int) ([]models.Student, int64, error) {
	// var student []models.Student
	var totalCount int64
	var args []any

	var whereBuilder strings.Builder
	whereBuilder.WriteString(" WHERE 1=1")

	student := make([]models.Student, 0, limit)

	if filters.FirstName != "" {
		args = append(args, filters.FirstName)
		fmt.Fprintf(&whereBuilder, " AND first_name = $%d", len(args))
	}
	if filters.LastName != "" {
		args = append(args, filters.LastName)
		fmt.Fprintf(&whereBuilder, " AND last_name = $%d", len(args))
	}
	if filters.Email != "" {
		args = append(args, filters.Email)
		fmt.Fprintf(&whereBuilder, " AND email = $%d", len(args))
	}
	slog.Debug("args", "args", args)

	countQuery := "SELECT COUNT(*) FROM students" + whereBuilder.String()
	err := p.pgx.QueryRow(ctx, countQuery, args...).Scan(&totalCount)
	if err != nil {
		slog.Debug("count query", "err", err)
		return nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	var sqlBuilder strings.Builder

	sqlBuilder.WriteString("SELECT id, first_name, last_name, email FROM students")
	sqlBuilder.WriteString(whereBuilder.String())

	if len(filters.SortBy) > 0 {
		sqlBuilder.WriteString(" ORDER BY ")
		for i, s := range filters.SortBy {
			if i > 0 {
				sqlBuilder.WriteString(", ")
			}
			fmt.Fprintf(&sqlBuilder, "%s %s", s.Field, s.Order)
		}
	}

	offset := (page - 1) * limit
	args = append(args, limit, offset)
	fmt.Fprintf(&sqlBuilder, " LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := p.pgx.Query(ctx, sqlBuilder.String(), args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query students: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Student
		err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email)
		if err != nil {
			return nil, 0, fmt.Errorf("scan student: %w", err)
		}
		// alloc_space from wrk is so big!
		// add cap in make()
		student = append(student, s)
	}
	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	// tx := p.db.WithContext(ctx).Model(&student)
	// if filters.FirstName != "" {
	// 	tx = tx.Where("first_name = ?", filters.FirstName)
	// }
	// if filters.LastName != "" {
	// 	tx = tx.Where("last_name = ?", filters.LastName)
	// }
	// if filters.Email != "" {
	// 	tx = tx.Where("email = ?", filters.Email)
	// }

	// if err := tx.Count(&totalCount).Error; err != nil {
	// 	return nil, 0, err
	// }

	// for _, sort := range filters.SortBy {
	// 	tx = tx.Order(sort.Field + " " + sort.Order)
	// }

	// result := tx.Limit(limit).Offset(offset).Find(&student)
	// if result.Error != nil {
	// 	return nil, 0, result.Error
	// }

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
	// result := p.db.Create(newStudent)
	// if result.Error != nil {
	// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 		return nil, result.Error
	// 	}
	// 	return nil, result.Error
	// }

	if err := p.db.Create(newStudent).Error; err != nil {
		return nil, fmt.Errorf("cannot add student: %w", err)
	}

	// 	if err := p.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
	// 	return 0, fmt.Errorf("get user: %w", err)
	// }

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
