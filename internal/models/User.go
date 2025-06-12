package models

import (
	"gamershub/internal/types"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id           uint              `gorm:"primaryKey" json:"id"`
	Username     string            `gorm:"unique;not null" json:"username"`
	Birthday     time.Time         `gorm:"type:timestamptz" json:"birthday,omitempty"`
	Email        types.Email       `gorm:"uniqueIndex;not null" json:"email"`
	PhoneNumber  types.PhoneNumber `gorm:"not null;" json:"phone_number"`
	Password     string            `json:"password"`
	Rating       float32           `gorm:"default: null" json:"rating"`
	Role         types.Role        `gorm:"type:varchar(20);default:'ROLE_USER'"`
	RefreshToken string            `gorm:"size:500"`
	TokenExpiry  time.Time
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest содержит данные для регистрации пользователя
type RegisterRequest struct {
	// Имя пользователя
	Username string `json:"username" example:"john_doe" binding:"required"`
	// Дата рождения в формате YYYY-MM-DD
	Birthday string `json:"birthday" example:"1990-01-01" binding:"required"`
	// Email пользователя (должен быть валидным email)
	Email types.Email `json:"email" example:"user@example.com" binding:"required"`
	// Пароль (минимум 6 символов)
	Password string `json:"password" example:"securePassword123" binding:"required,min=6"`
	// Номер телефона (должен быть валидным номером)
	PhoneNumber types.PhoneNumber `json:"phone_number" example:"+1234567890" binding:"required"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required,min=8,max=64"`
	NewPassword     string `json:"newPassword" binding:"required,min=8,max=64,nefield=CurrentPassword"`
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type RegisterSuccessResponse struct {
	// Статус операции
	Status string `json:"status" example:"success"`
	// Сообщение
	Message string `json:"message" example:"User registered successfully"`
	// Данные пользователя
	Data UserDataResponse `json:"data"`
}

// UserDataResponse содержит возвращаемые сервером данные регистрации
type UserDataResponse struct {
	UserID      uint              `json:"user_id"`
	Email       types.Email       `json:"email"`
	Username    string            `json:"username"`
	Birthday    string            `json:"birthday"`
	PhoneNumber types.PhoneNumber `json:"phoneNumber"`
	Token       string            `json:"token"`
}

// ErrorResponse стандартный ответ с ошибкой
type ErrorResponse struct {
	// Сообщение об ошибке
	Error string `json:"error" example:"Invalid email format"`
}
