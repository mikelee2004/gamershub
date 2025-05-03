package controllers

import (
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gamershub/internal/types"
	"gamershub/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	userRepository *repositories.UserRepository
}

func NewAuthController(userRepo *repositories.UserRepository) *AuthController {
	return &AuthController{userRepository: userRepo}
}

// Login method
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

	// Остальная логика...
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// creating validated email
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
		Email:    email,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     types.RoleUser,
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
			"user_id": user.ID,
			"email":   user.Email.String(),
			"token":   token,
		},
	})
}
