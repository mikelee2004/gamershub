package models

import "gorm.io/gorm"

type GamingSession struct {
	gorm.Model
	InitiatorId   uint    `json:"initiator_id" binding:"required"`
	PartnerId     uint    `json:"partner_id"`
	MatchResult   string  `gorm:"type:text"`
	InitiatorRate float32 `json:"initiator_rate"`
	PartnerRate   float32 `json:"partner_rate"`
	Status        string  `json:"status"`
}
