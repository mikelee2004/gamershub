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

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя в системе
// @Tags Auth
// @Accept json
// @Produce json
// @security none
// @Param registerRequest body models.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} models.RegisterSuccessResponse "Успешная регистрация"
// @Failure 400 {object} models.ErrorResponse "Невалидные данные"
// @Failure 409 {object} models.ErrorResponse "Пользователь уже существует"
// @Failure 500 {object} models.ErrorResponse "Внутренняя ошибка сервера"
// @Router /auth/register [post]
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

	phone, err := types.NewPhoneNumber(string(req.PhoneNumber))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	email, err := types.NewEmail(string(req.Email))
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
	res := models.RegisterSuccessResponse{
		Status:  "success",
		Message: "User registered successfully",
		Data: models.UserDataResponse{
			UserID:      user.Id,
			Email:       user.Email,
			Username:    user.Username,
			Birthday:    user.Birthday.String(),
			PhoneNumber: user.PhoneNumber,
			Token:       token,
		},
	}
	c.JSON(http.StatusCreated, res)
}

// Logout godoc
// @Summary User logout
// @Description Invalidates the user's refresh token and logs them out
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.LogoutSuccessResponse
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/logout [post]
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

// Login godoc
// @Summary User authentication
// @Description Authenticate user with email and password. Returns JWT tokens.
// @Tags Auth
// @Accept json
// @Produce json
// @security none
// @Param request body models.EmailLoginRequest true "Login credentials"
// @Success 200 {object} models.LoginSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.EmailLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//email validation
	email, err := types.NewEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// checking validated email
	user, err := ac.userRepository.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// generating tokens
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

type RegisterSuccessResponse struct {
	// Статус операции
	Status string `json:"status" example:"success"`
	// Сообщение
	Message string `json:"message" example:"User registered successfully"`
	// Данные пользователя
	Data struct {
		// ID пользователя
		UserID string `json:"user_id" example:"507f1f77bcf86cd799439011"`
		// Email пользователя
		Email string `json:"email" example:"user@example.com"`
		// Имя пользователя
		Username string `json:"username" example:"john_doe"`
		// Дата рождения
		Birthday string `json:"birthday" example:"1990-01-01"`
		// Номер телефона
		PhoneNumber string `json:"phoneNumber" example:"+1234567890"`
		// JWT токен
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	} `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid email format"`
}
