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
	Email    string `json:"email" binding:"required,email" example:"mike@mail.ru"`
	Password string `json:"password" binding:"required" example:"123123"`
}

// RegisterRequest содержит данные для регистрации пользователя
type RegisterRequest struct {
	Username    string            `json:"username" example:"test1" binding:"required"`
	Birthday    string            `json:"birthday" example:"1990-01-01" binding:"required"`
	Email       types.Email       `json:"email" example:"user@test.com" binding:"required"`
	Password    string            `json:"password" example:"123123" binding:"required,min=6"`
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

type LoginSuccessResponse struct {
	// Статус операции
	Status string `json:"status" example:"success"`
	// Сообщение
	Message string `json:"message" example:"User registered successfully"`
	// Данные пользователя
	Data LoginDataResponse `json:"data"`
}

type LoginDataResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
