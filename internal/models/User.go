package models

import (
	"gamershub/internal/types"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Id          uint              `gorm:"primaryKey" json:"id"`
	Username    string            `gorm:"unique;not null" json:"username"`
	Birthday    time.Time         `gorm:"type:timestamptz" json:"birthday,omitempty"`
	Email       types.Email       `gorm:"uniqueIndex;not null" json:"email"`
	PhoneNumber types.PhoneNumber `gorm:"not null;" json:"phone_number"`
	Password    string            `json:"password"`
	Rating      float32           `gorm:"not null" json:"rating"`
	Role        types.Role        `gorm:"type:varchar(20);default:'user'"`
}

type EmailLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Birthday    string `json:"birthday" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required,min=8"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}
