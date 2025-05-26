package service

import (
	"testing"

	"github.com/one2n/student-api/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func TestCreateStudent(t *testing.T) {
	db := setupTestDB(t)
	service := NewStudentService(db)

	tests := []struct {
		name    string
		student *model.Student
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid student",
			student: &model.Student{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   20,
				Grade: "A",
			},
			wantErr: false,
		},
		{
			name: "invalid age",
			student: &model.Student{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   0, // Below minimum age
				Grade: "A",
			},
			wantErr: true,
			errMsg:  "invalid age: must be between 1 and 100",
		},
		{
			name: "duplicate email",
			student: &model.Student{
				Name:  "Jane Doe",
				Email: "john@example.com", // Same email as first test
				Age:   21,
				Grade: "B",
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			student, err := service.CreateStudent(tt.student)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, student.ID)
			assert.Equal(t, tt.student.Name, student.Name)
			assert.Equal(t, tt.student.Email, student.Email)
			assert.Equal(t, tt.student.Age, student.Age)
			assert.Equal(t, tt.student.Grade, student.Grade)
		})
	}
}

func TestGetStudentByID(t *testing.T) {
	db := setupTestDB(t)
	service := NewStudentService(db)

	// Create a test student
	student := &model.Student{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   20,
		Grade: "A",
	}
	createdStudent, err := service.CreateStudent(student)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing student",
			id:      createdStudent.ID,
			wantErr: false,
		},
		{
			name:    "non-existing student",
			id:      "non-existing-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			student, err := service.GetStudentByID(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.id, student.ID)
		})
	}
}

func TestGetAllStudents(t *testing.T) {
	db := setupTestDB(t)
	service := NewStudentService(db)

	// Create test students
	students := []*model.Student{
		{Name: "John Doe", Email: "john@example.com", Age: 20, Grade: "A"},
		{Name: "Jane Smith", Email: "jane@example.com", Age: 21, Grade: "B"},
	}

	for _, student := range students {
		_, err := service.CreateStudent(student)
		assert.NoError(t, err)
	}

	// Test getting all students
	retrievedStudents, err := service.GetAllStudents()
	assert.NoError(t, err)
	assert.Len(t, retrievedStudents, len(students))
}

func TestUpdateStudent(t *testing.T) {
	db := setupTestDB(t)
	service := NewStudentService(db)

	// Create test students
	student1 := &model.Student{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   20,
		Grade: "A",
	}
	student2 := &model.Student{
		Name:  "Jane Smith",
		Email: "jane@example.com",
		Age:   21,
		Grade: "B",
	}
	createdStudent1, err := service.CreateStudent(student1)
	assert.NoError(t, err)
	_, err = service.CreateStudent(student2)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		update  *model.Student
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update",
			id:   createdStudent1.ID,
			update: &model.Student{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Age:   21,
				Grade: "A+",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			id:   createdStudent1.ID,
			update: &model.Student{
				Name:  "John Updated",
				Email: "jane@example.com", // Using student2's email
				Age:   21,
				Grade: "A+",
			},
			wantErr: true,
			errMsg:  "email already exists",
		},
		{
			name: "non-existing student",
			id:   "non-existing-id",
			update: &model.Student{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Age:   21,
				Grade: "A+",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updatedStudent, err := service.UpdateStudent(tt.id, tt.update)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.EqualError(t, err, tt.errMsg)
				}
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.update.Name, updatedStudent.Name)
			assert.Equal(t, tt.update.Email, updatedStudent.Email)
			assert.Equal(t, tt.update.Age, updatedStudent.Age)
			assert.Equal(t, tt.update.Grade, updatedStudent.Grade)
		})
	}
}

func TestDeleteStudent(t *testing.T) {
	db := setupTestDB(t)
	service := NewStudentService(db)

	// Create a test student
	student := &model.Student{
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   20,
		Grade: "A",
	}
	createdStudent, err := service.CreateStudent(student)
	assert.NoError(t, err)

	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{
			name:    "existing student",
			id:      createdStudent.ID,
			wantErr: false,
		},
		{
			name:    "non-existing student",
			id:      "non-existing-id",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.DeleteStudent(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify student is deleted
			_, err = service.GetStudentByID(tt.id)
			assert.Error(t, err)
		})
	}
}
