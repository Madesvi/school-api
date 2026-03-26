package students

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	// ЗАПУСК
	// go test -v ./internal/api/handlers/students

	// Путь к вашим новым мокам
	"rest-api-app/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetOneStudentHandler(t *testing.T) {
	// Описываем таблицу тестов
	tests := []struct {
		name           string
		studentID      string
		mockID         int
		mockReturn     models.Student
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Success",
			studentID:      "1",
			mockID:         1,
			mockReturn:     models.Student{ID: 1, FirstName: "Ivan"},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Not Found",
			studentID:      "999",
			mockID:         999,
			mockReturn:     models.Student{},
			mockErr:        errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Invalid ID format",
			studentID:      "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockGetStudentByID(t)

			// Настраиваем мок только если ID валидный (не "abc")
			if tt.studentID != "abc" {
				mockService.On("GetStudentByID", mock.Anything, tt.mockID).
					Return(tt.mockReturn, tt.mockErr)
			}

			handler := GetOneStudentHandler(mockService)
			req := httptest.NewRequest(http.MethodGet, "/students/"+tt.studentID, nil)
			req.SetPathValue("id", tt.studentID)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}

// // ТЕСТ 1: Проверяем успешный сценарий
// func TestGetOneStudentHandler_Success(t *testing.T) {
// 	mockService := NewMockGetStudentByID(t)
//
// 	expectedStudent := models.Student{ID: 1, FirstName: "Ivan"}
//
// 	// Ожидаем ID = 1
// 	mockService.On("GetStudentByID", mock.Anything, 1).
// 		Return(expectedStudent, nil)
//
// 	handler := GetOneStudentHandler(mockService)
//
// 	req := httptest.NewRequest(http.MethodGet, "/students/1", nil)
// 	req.SetPathValue("id", "1") // Здесь ID должен совпадать с моком!
// 	rr := httptest.NewRecorder()
//
// 	handler.ServeHTTP(rr, req)
//
// 	assert.Equal(t, http.StatusOK, rr.Code)
// }
//
// // ТЕСТ 2: Проверяем сценарий "Студент не найден"
// func TestGetOneStudentHandler_NotFound(t *testing.T) {
// 	mockService := NewMockGetStudentByID(t)
//
// 	// Ожидаем ID = 999 и возвращаем ошибку
// 	mockService.On("GetStudentByID", mock.Anything, 999).
// 		Return(models.Student{}, errors.New("not found"))
//
// 	handler := GetOneStudentHandler(mockService)
//
// 	req := httptest.NewRequest(http.MethodGet, "/students/999", nil)
// 	req.SetPathValue("id", "999") // Здесь тоже 999!
// 	rr := httptest.NewRecorder()
//
// 	handler.ServeHTTP(rr, req)
//
// 	// Здесь проверяем, что статус именно 404 (если вы добавили http.Error в хэндлер)
// 	assert.Equal(t, http.StatusNotFound, rr.Code)
// }
//
// func TestGetOneStudentHandler_InvalidID(t *testing.T) {
// 	mockService := NewMockGetStudentByID(t) // Мок создаем, но вызывать не будем
// 	handler := GetOneStudentHandler(mockService)
//
// 	// Передаем "abc" вместо цифры
// 	req := httptest.NewRequest(http.MethodGet, "/students/abc", nil)
// 	req.SetPathValue("id", "abc")
// 	rr := httptest.NewRecorder()
//
// 	handler.ServeHTTP(rr, req)
//
// 	// Ожидаем 400 Bad Request
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// }

// func TestGetOneStudentHandler_Success(t *testing.T) {
// 	// 1. Создаем мок из сгенерированного файла
// 	mockService := NewMockGetStudentByID(t)
//
// 	// 2. Настраиваем ожидание (Expectation)
// 	expectedStudent := models.Student{ID: 1, FirstName: "Ivan"}
//
// 	// Говорим моку: "Когда тебя вызовут с любым контекстом и ID=1, верни студента и nil"
// 	mockService.On("GetStudentByID", mock.Anything, 1).
// 		Return(expectedStudent, nil)
//
// 	// 3. Создаем хэндлер и передаем ему наш мок
// 	handler := GetOneStudentHandler(mockService)
//
// 	// 4. Подготавливаем HTTP запрос (Go 1.22+)
// 	req := httptest.NewRequest(http.MethodGet, "/students/1", nil)
// 	req.SetPathValue("id", "1") // Эмулируем параметр пути
//
// 	rr := httptest.NewRecorder()
//
// 	// 5. Запускаем
// 	handler.ServeHTTP(rr, req)
//
// 	// 6. Проверяем результат
// 	assert.Equal(t, http.StatusOK, rr.Code)
//
// 	var actualStudent models.Student
// 	json.NewDecoder(rr.Body).Decode(&actualStudent)
// 	assert.Equal(t, expectedStudent.FirstName, actualStudent.FirstName)
//
// 	// Проверяем, что метод мока действительно вызывался
// 	mockService.AssertExpectations(t)
// }
