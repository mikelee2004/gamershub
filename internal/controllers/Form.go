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

// CreateForm godoc
// @Summary User authentication
// @Description Authenticate user with email and password. Returns JWT tokens.
// @Tags Player Forms
// @Accept json
// @Produce json
// @Param request body models.PlayerForm true "Form Credentials"
// @Success 200 {object} models.FormCreatedResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /form/create [post]
func (c *PlayerFormController) CreateForm(ctx *gin.Context) {
	var form models.PlayerForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}
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
	ctx.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   form,
	})
}

// Get godoc
// @Summary Get player form
// @Description Retrieves the current player's form with all preferences and statistics
// @Tags Player Forms
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.PlayerForm "Player form data"
// @Failure 404 {object} models.ErrorResponse "Form not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /form/me [get]
func (c *PlayerFormController) Get(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	form, err := c.repo.GetByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Form not found"})
		return
	}

	ctx.JSON(http.StatusOK, form)
}

// Delete godoc
// @Summary Delete player form
// @Description Permanently deletes the player's form
// @Tags Player Forms
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.DeleteSuccessResponse "Success message"
// @Failure 404 {object} models.ErrorResponse "Form not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /form/delete [delete]
func (c *PlayerFormController) Delete(ctx *gin.Context) {
	userID := ctx.GetUint("userID")
	if err := c.repo.Delete(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete form"})
		return
	}
	ctx.JSON(http.StatusOK, models.DeleteSuccessResponse{
		Status:  "success",
		Message: "Анкета удалена",
	})
}
