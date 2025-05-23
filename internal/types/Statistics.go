package types

import "gorm.io/gorm"

type Statistics struct {
	gorm.Model
	UserId    uint    `gorm:"not null;uniqueIndex" json:"player_id"`
	Rank      string  `gorm:"size:20" json:"rank" validate:"required"`
	PeakRank  string  `gorm:"size:20" json:"peakRank"`
	KD        float32 `gorm:"type:numeric(4,2)" json:"kd" validate:"gte=0"`
	WinRate   float32 `gorm:"type:numeric(4,2)" json:"win_rate" validate:"gte=0,lte=100"`
	HSPercent float32 `gorm:"type:numeric(4,2)" json:"hs_percent" validate:"gte=0,lte=100"`
}
