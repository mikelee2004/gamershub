package models

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type PlayerForm struct {
	UserID            uint       `gorm:"not null;uniqueIndex" json:"user_id"`
	ValorantID        string     `gorm:"not null" json:"valorant_id" example:"Testing#gogin"`
	Goal              string     `gorm:"type:varchar(50)" json:"goal" example:"ranked/pro/chill"`
	StatisticsJSON    string     `gorm:"type:jsonb" json:"-"`
	CommunicationJSON string     `gorm:"type:jsonb" json:"-"`
	RegionJSON        string     `gorm:"type:jsonb" json:"-"`
	Statistics        Statistics `gorm:"-" json:"statistics"`
	CommunicationWays []string   `gorm:"-" json:"communication_ways"`
	Region            Region     `gorm:"-" json:"region"`
	Description       string     `gorm:"type:text" json:"description"`
}

func (p *PlayerForm) DeserializeJSONData() error {
	// Убираем лишние экранирования если они есть
	log.Printf("Deserializing for user %d:", p.UserID)
	log.Printf("Raw StatisticsJSON: %s", p.StatisticsJSON)
	log.Printf("Raw CommunicationJSON: %s", p.CommunicationJSON)
	log.Printf("Raw RegionJSON: %s", p.RegionJSON)
	statsJSON := cleanEscapedJSON(p.StatisticsJSON)
	commJSON := cleanEscapedJSON(p.CommunicationJSON)
	regionJSON := cleanEscapedJSON(p.RegionJSON)

	// Десериализация Statistics
	if statsJSON != "" {
		if err := json.Unmarshal([]byte(statsJSON), &p.Statistics); err != nil {
			return fmt.Errorf("failed to unmarshal Statistics: %v, data: %s", err, statsJSON)
		}
	}

	// Десериализация CommunicationWays (специальная обработка)
	if commJSON != "" {
		// Убираем лишние кавычки если они есть
		commJSON = strings.ReplaceAll(commJSON, `\"`, `"`)

		// Обработка случая когда массив представлен как строка
		if strings.HasPrefix(commJSON, `"`) && strings.HasSuffix(commJSON, `"`) {
			commJSON = commJSON[1 : len(commJSON)-1]
		}

		if err := json.Unmarshal([]byte(commJSON), &p.CommunicationWays); err != nil {
			return fmt.Errorf("failed to unmarshal CommunicationWays: %v, data: %s", err, commJSON)
		}
	}

	// Десериализация Region
	if regionJSON != "" {
		if err := json.Unmarshal([]byte(regionJSON), &p.Region); err != nil {
			return fmt.Errorf("failed to unmarshal Region: %v, data: %s", err, regionJSON)
		}
	}
	log.Printf("Result for user %d: Rank=%d, Region=%s, Languages=%v, Comms=%v",
		p.UserID, p.Statistics.Rank, p.Region.Region,
		p.Region.Language, p.CommunicationWays)

	return nil
}

func cleanEscapedJSON(s string) string {
	// Убираем лишние экранирования если они есть
	if strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		s = s[1 : len(s)-1]
	}
	return strings.ReplaceAll(s, `\"`, `"`)
}

type Region struct {
	Region   string   `json:"region" example:"EU/NA/ASIA/AU"`
	Language []string `json:"language" example:"en/ru"`
}

// swagger

// FormCreatedResponse represents successful form creation response
type FormCreatedResponse struct {
	Status  string            `json:"status" example:"success"`
	Data    PlayerFormDetails `json:"data"`
	Message string            `json:"message,omitempty" example:"Form created successfully"`
}

// PlayerFormDetails contains subset of player form fields for response
type PlayerFormDetails struct {
	ID         uint   `json:"id" example:"123"`
	UserID     uint   `json:"user_id" example:"456"`
	ValorantID string `json:"valorant_id" example:"Player#1234"`
	Goal       string `json:"goal" example:"Improve skills"`
	CreatedAt  string `json:"created_at" example:"2023-05-15T14:23:01Z"`
}

type DeleteSuccessResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Resource deleted successfully"`
}
