package controllers

import (
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PlayerFormController struct {
	repo *repositories.FormRepository
}

func NewPlayerFormController(repo *repositories.FormRepository) *PlayerFormController {
	return &PlayerFormController{repo: repo}
}

func (c *PlayerFormController) CreateForm(ctx *gin.Context) {
	var form models.PlayerForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	// Получаем userID из контекста (JWT middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	form.UserID = userID.(uint)

	if err := c.repo.Create(&form); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create form",
			"details": err.Error(), // Возвращаем детали ошибки для отладки
		})
		return
	}

	// Возвращаем созданную форму с десериализованными данными
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   form,
	})
}

// Get возвращает анкету пользователя
func (c *PlayerFormController) Get(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	form, err := c.repo.GetByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	ctx.JSON(http.StatusOK, form)
}

// Delete удаляет анкету
func (c *PlayerFormController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint("userID")

	if err := c.repo.Delete(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
