package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/one2n/student-api/model"
	"gorm.io/gorm"
)

type StudentService struct {
	db *gorm.DB
}

func NewStudentService(db *gorm.DB) *StudentService {
	err := db.AutoMigrate(&model.Student{})
	if err != nil {
		fmt.Printf("Error migrating schema: %v\n", err)
	}
	return &StudentService{db: db}
}

func (s *StudentService) CreateStudent(student *model.Student) (*model.Student, error) {
	if student.Age <= 0 || student.Age > 100 {
		return nil, errors.New("invalid age: must be between 1 and 100")
	}

	// Check if email already exists
	var existingStudent model.Student
	if err := s.db.Where("email = ?", student.Email).First(&existingStudent).Error; err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("error checking email uniqueness: %v", err)
	}

	student.ID = uuid.New().String()
	student.CreatedAt = time.Now()
	student.UpdatedAt = time.Now()

	result := s.db.Create(student)
	if result.Error != nil {
		return nil, result.Error
	}
	return student, nil
}

func (s *StudentService) GetAllStudents() ([]*model.Student, error) {
	var students []*model.Student
	result := s.db.Find(&students)
	if result.Error != nil {
		return nil, result.Error
	}
	return students, nil
}

func (s *StudentService) GetStudentByID(id string) (*model.Student, error) {
	var student model.Student
	result := s.db.First(&student, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, result.Error
	}
	return &student, nil
}

func (s *StudentService) UpdateStudent(id string, updatedStudent *model.Student) (*model.Student, error) {
	var student model.Student
	result := s.db.First(&student, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, result.Error
	}

	if updatedStudent.Email != student.Email {
		var existingStudent model.Student
		if err := s.db.Where("email = ? AND id != ?", updatedStudent.Email, id).First(&existingStudent).Error; err == nil {
			return nil, errors.New("email already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("error checking email uniqueness: %v", err)
		}
	}

	student.Name = updatedStudent.Name
	student.Email = updatedStudent.Email
	student.Age = updatedStudent.Age
	student.Grade = updatedStudent.Grade
	student.UpdatedAt = time.Now()

	result = s.db.Save(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}

func (s *StudentService) DeleteStudent(id string) error {
	result := s.db.Delete(&model.Student{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("student not found")
	}
	return nil
}
