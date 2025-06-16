package models

import (
	"gorm.io/gorm"
	"time"
)

type Match struct {
	gorm.Model
	GameSessionID uint      `gorm:"not null;uniqueIndex" json:"game_session_id"`
	Player1ID     uint      `gorm:"not null" json:"player1_id"`
	Player2ID     uint      `gorm:"not null" json:"player2_id"`
	WinnerID      *uint     `json:"winner_id"`                      // nil для ничьи
	Result        string    `gorm:"type:varchar(10)" json:"result"` // "win", "loss", "draw"
	Score         string    `gorm:"type:varchar(20)" json:"score"`  // "13-5"
	Duration      int       `gorm:"not null" json:"duration"`       // в минутах
	EndedAt       time.Time `gorm:"not null" json:"ended_at"`
	Stats         string    `gorm:"type:json" json:"stats`
	Player1       User      `gorm:"foreignKey:Player1ID" json:"-"`
	Player2       User      `gorm:"foreignKey:Player2ID" json:"-"`
	Winner        *User     `gorm:"foreignKey:WinnerID" json:"-"`
}
