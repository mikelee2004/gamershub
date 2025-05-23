package models

import (
	"gorm.io/gorm"
)

type PlayerForm struct {
	gorm.Model
	UserID     uint   `gorm:"not null;uniqueIndex" json:"user_id"`
	ValorantID string `gorm:"not null" json:"valorant_id"`
	Goal       string `gorm:"type:varchar(50)" json:"goal"`

	// Храним как JSON в БД
	StatisticsJSON    string `gorm:"type:jsonb" json:"-"`
	CommunicationJSON string `gorm:"type:jsonb" json:"-"`
	RegionJSON        string `gorm:"type:jsonb" json:"-"`

	// Для API (игнорируем в GORM)
	Statistics        Statistics `gorm:"-" json:"statistics"`
	CommunicationWays []string   `gorm:"-" json:"communication_ways"`
	Region            Region     `gorm:"-" json:"region"`
	Description       string     `gorm:"type:text" json:"description"`
}

type Region struct {
	Region   string   `json:"region"`
	Timezone string   `json:"timezone"`
	Language []string `json:"language"`
}
