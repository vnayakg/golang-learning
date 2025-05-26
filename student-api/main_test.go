package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/one2n/student-api/model"
	"github.com/one2n/student-api/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter() (*gin.Engine, *service.StudentService) {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	studentService := service.NewStudentService(db)
	r := setupRouter(studentService)
	return r, studentService
}

func TestCreateStudentHandler(t *testing.T) {
	r, _ := setupTestRouter()

	tests := []struct {
		name       string
		payload    model.Student
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid student",
			payload: model.Student{
				Name:  "John Doe",
				Email: "john@doe.com",
				Age:   20,
				Grade: "A",
			},
			wantStatus: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "invalid age",
			payload: model.Student{
				Name:  "John Doe",
				Email: "john@doe.com",
				Age:   5,
				Grade: "A",
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/v1/api/students", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response model.StudentResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotNil(t, response.Data)
			}
		})
	}
}

func TestGetAllStudentsHandler(t *testing.T) {
	r, service := setupTestRouter()

	// Create test students with unique emails
	students := []*model.Student{
		{Name: "John Doe", Email: "john@doe.com", Age: 20, Grade: "A"},
		{Name: "Jane Smith", Email: "jane@doe.com", Age: 21, Grade: "B"},
	}

	for _, student := range students {
		_, err := service.CreateStudent(student)
		assert.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/api/students", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.StudentResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
}

func TestGetStudentByIDHandler(t *testing.T) {
	r, service := setupTestRouter()

	student := &model.Student{
		Name:  "John Doe",
		Email: "john@doe.com",
		Age:   20,
		Grade: "A",
	}
	createdStudent, err := service.CreateStudent(student)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "existing student",
			id:         createdStudent.ID,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing student",
			id:         "non-existing-id",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/v1/api/students/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response model.StudentResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotNil(t, response.Data)
			}
		})
	}
}

func TestUpdateStudentHandler(t *testing.T) {
	r, service := setupTestRouter()

	// Create a test student
	student := &model.Student{
		Name:  "John Doe",
		Email: "john@doe.com",
		Age:   20,
		Grade: "A",
	}
	createdStudent, err := service.CreateStudent(student)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		payload    model.Student
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid update",
			id:   createdStudent.ID,
			payload: model.Student{
				Name:  "John Updated",
				Email: "john.updated@doe.com",
				Age:   21,
				Grade: "A+",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "non-existing student",
			id:   "non-existing-id",
			payload: model.Student{
				Name:  "John Updated",
				Email: "john.updated@doe.com",
				Age:   21,
				Grade: "A+",
			},
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPut, "/v1/api/students/"+tt.id, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response model.StudentResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
				assert.NotNil(t, response.Data)
			}
		})
	}
}

func TestDeleteStudentHandler(t *testing.T) {
	r, service := setupTestRouter()

	student := &model.Student{
		Name:  "John Doe",
		Email: "john@doe.com",
		Age:   20,
		Grade: "A",
	}
	createdStudent, err := service.CreateStudent(student)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		id         string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "existing student",
			id:         createdStudent.ID,
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name:       "non-existing student",
			id:         "non-existing-id",
			wantStatus: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/v1/api/students/"+tt.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
			if !tt.wantErr {
				var response model.StudentResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.True(t, response.Success)
			}
		})
	}
}
