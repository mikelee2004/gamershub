package models

import "gorm.io/gorm"

type Statistics struct {
	Rank      uint    `json:"rank"`
	PeakRank  uint    `json:"peak_rank"`
	KD        float32 `json:"kd"`
	WinRate   float32 `json:"win_rate"`
	HSPercent float32 `json:"hs_percent"`
}

type Rank struct {
	gorm.Model
	Name      string `gorm:"not null" json:"rank_name"`
	Division  int    `gorm:"not null" json:"rank_division"`
	RankValue int    `gorm:"not null" json:"rank_value"`
}

type PlayerTrait struct {
	gorm.Model
	PlayerID        uint    `gorm:"not null;uniqueIndex"`
	StressTolerance float64 `gorm:"default:5.0"` // 1-10
	Talkativeness   float64 `gorm:"default:5.0"` // 1-10
	PreferredGoal   string  // "ranked", "fun", "training/pro"
	PreferredRoles  string  // "duelist,senti,initiator,controller"
}

type LogoutSuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Successfully logged out"`
}
