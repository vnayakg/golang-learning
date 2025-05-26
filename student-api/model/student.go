package model

import (
	"time"

	"gorm.io/gorm"
)

type Student struct {
	ID        string         `json:"id" gorm:"primaryKey;type:text"`
	Name      string         `json:"name" gorm:"not null" binding:"required"`
	Email     string         `json:"email" gorm:"not null;uniqueIndex" binding:"required,email"`
	Age       int            `json:"age" gorm:"not null" binding:"required,min=6"`
	Grade     string         `json:"grade" gorm:"not null" binding:"required"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type StudentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}
