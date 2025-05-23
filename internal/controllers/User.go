package controllers

import (
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gamershub/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type UserController struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB
}

func NewUserController(userRepo *repositories.UserRepository, db *gorm.DB) *UserController {
	return &UserController{userRepo: userRepo, db: db}
}

func (u *UserController) GetProfile(c *gin.Context) {
	userId := c.GetUint("userID")
	user, err := u.userRepo.FindUserByID(userId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"email":    user.Email,
	})
}

func (u *UserController) ChangePassword(c *gin.Context) {
	userID := c.GetUint("userID")
	var req models.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := u.db.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	// Проверка текущего пароля
	if err := utils.ComparePasswords(user.Password, req.CurrentPassword); err != nil {
		c.JSON(401, gin.H{"error": "Invalid current password"})
		return
	}

	// Хеширование нового пароля
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to process password"})
		return
	}

	// Обновление пароля в транзакции
	tx := u.db.Begin()
	if err := tx.Model(&user).Update("password", hashedPassword).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	// Инвалидация всех активных токенов (опционально)
	if err := tx.Model(&user).Update("token_expiry", time.Now()).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to invalidate sessions"})
		return
	}

	tx.Commit()

	c.JSON(200, models.ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}
