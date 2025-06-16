package controllers

import (
	"gamershub/internal/repositories"
	"gorm.io/gorm"
)

type GamingSessionController struct {
	playerForm *repositories.FormRepository
	db         *gorm.DB
}

func NewGamingSessionController(playerForm *repositories.FormRepository, db *gorm.DB) {

}
