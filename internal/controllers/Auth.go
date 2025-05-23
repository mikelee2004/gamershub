package controllers

import (
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gamershub/internal/types"
	"gamershub/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type AuthController struct {
	userRepository *repositories.UserRepository
}

func NewAuthController(userRepo *repositories.UserRepository) *AuthController {
	return &AuthController{userRepository: userRepo}
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email, err := types.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Теперь передаем корректный Email
	user, err := ac.userRepository.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	accessToken, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации токена"})
	}
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации токена"})
		return
	}
	user.RefreshToken = refreshToken
	user.TokenExpiry = time.Now().Add(7 * 24 * time.Hour)

	if err := ac.userRepository.UpdateUser(user); err != nil {
		log.Printf("UpdateUser error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения токена"})
		return
	}

	updatedUser, err := ac.userRepository.FindUserByID(user.Id)
	if err != nil {
		log.Printf("Failed to verify saved token: %v", err)
	} else {
		log.Printf("Refresh token in DB after save: %s", updatedUser.RefreshToken)
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	birthday, err := types.ValidateBirthday(req.Birthday)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	phone, err := types.NewPhoneNumber(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	email, err := types.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}

	user := &models.User{
		Email:       email,
		Username:    req.Username,
		Birthday:    birthday,
		Password:    string(hashedPassword),
		PhoneNumber: phone,
		Role:        types.RoleUser,
	}

	if err := ac.userRepository.CreateUser(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User registered successfully",
		"data": gin.H{
			"user_id":     user.Id,
			"email":       user.Email.String(),
			"username":    user.Username,
			"birthday":    user.Birthday,
			"phoneNumber": user.PhoneNumber,
			"token":       token,
		},
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	userId := c.GetUint("userID")
	if err := ac.userRepository.RevokeRefreshToken(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not revoke refresh token, logout failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out",
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No refresh token found"})
		return
	}
	user, err := ac.userRepository.FindByRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if time.Now().After(user.TokenExpiry) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
		return
	}
	newAccessToken, _ := utils.GenerateJWT(user)
	newRefreshToken, _ := utils.GenerateRefreshToken()

	user.RefreshToken = newRefreshToken

	if err := ac.userRepository.UpdateUser(user); err != nil {
		log.Printf("UpdateUser error: %v", err) // Логируем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сохранения токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
