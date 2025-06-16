package models

import (
	"gorm.io/gorm"
	"time"
)

type GameSession struct {
	gorm.Model
	SenderID    uint       `gorm:"not null" json:"sender_id"`
	RecipientID uint       `gorm:"not null" json:"recipient_id" example:"3"`
	Status      string     `gorm:"type:varchar(20);default:'pending';check:status IN ('pending','accepted','declined','expired','canceled')" json:"status"`
	Message     string     `gorm:"type:text" json:"message" example:"Бро, давай вместе поднимать ранг!"`
	ExpiresAt   time.Time  `gorm:"not null" json:"expires_at"`
	AcceptedAt  *time.Time `json:"accepted_at"`
	StartedAt   *time.Time `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at"`
}

type MatchCriteria struct {
	PlayerForm
	MaxDistance int
}

type MatchResult struct {
	Player     PlayerForm
	MatchScore float64
}

// swagger models

type InvitationResponse struct {
	InvitationID uint       `json:"invitation_id"`
	Sender       SenderInfo `json:"sender"`
	Message      string     `json:"message"`
	SentAt       time.Time  `json:"sent_at"`
	ExpiresAt    time.Time  `json:"expires_at"`
}

type SenderInfo struct {
	UserID     uint   `json:"user_id"`
	ValorantID string `json:"valorant_id"`
	Rank       uint   `json:"rank"`
	RankName   string `json:"rank_name"`
}

type RespondRequest struct {
	InvitationID uint   `json:"invitation_id" binding:"required"`
	Action       string `json:"action" binding:"required,oneof=accept decline"`
}

type CompleteRequest struct {
	SessionID uint   `json:"session_id" binding:"required"`
	WinnerID  *uint  `json:"winner_id"`
	Score     string `json:"score" binding:"required"`
	Stats     string `json:"stats"`
}

type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Bad request"`
}

type RespondInvitationRequest struct {
	// ID приглашения
	InvitationID uint `json:"invitation_id" binding:"required" example:"123"`

	// Действие: accept или decline
	Action string `json:"action" binding:"required,oneof=accept decline" example:"accept"`

	// Причина отказа (опционально)
	Reason *string `json:"reason,omitempty" example:"Not enough time"`
}

// MatchResponse represents successful match response
type MatchResponse struct {
	Count   int             `json:"count" example:"5"`
	Matches []MatchedPlayer `json:"matches"`
}

// MatchedPlayer represents matched player data
type MatchedPlayer struct {
	UserID     uint    `json:"user_id" example:"123"`
	ValorantID string  `json:"valorant_id" example:"Player#EUW"`
	Rank       uint    `json:"rank" example:"25"`
	KD         float64 `json:"kd" example:"1.25"`
	MatchScore float64 `json:"match_score" example:"0.85"`
}

// NoMatchesResponse represents response when no matches found
type NoMatchesResponse struct {
	Message      string         `json:"message" example:"No matches found"`
	YourCriteria PlayerCriteria `json:"your_criteria"`
}

// PlayerCriteria represents search criteria used
type PlayerCriteria struct {
	Region    string   `json:"region" example:"Europe"`
	Languages []string `json:"languages" example:"['English','Russian']"`
	Rank      uint     `json:"rank" example:"25"`
	Goal      string   `json:"goal" example:"Competitive"`
}

// SendInviteRequest represents invitation request payload
type SendInviteRequest struct {
	RecipientID uint   `json:"recipient_id" binding:"required" example:"123"`
	Message     string `json:"message" example:"Let's play together!" maxLength:"500"`
}

// InvitationSentResponse represents successful invitation response
type InvitationSentResponse struct {
	Message       string             `json:"message" example:"Invitation sent successfully"`
	InvitationID  uint               `json:"invitation_id" example:"1"`
	ExpiresAt     time.Time          `json:"expires_at" example:"2023-05-17T14:23:01Z"`
	RecipientInfo RecipientShortInfo `json:"recipient_info"`
}

type RecipientShortInfo struct {
	ID         uint   `json:"id" example:"123"`
	ValorantID string `json:"valorant_id" example:"Player#EUW"`
	Rank       uint   `json:"rank" example:"25"`
}
