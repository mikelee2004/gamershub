package controllers

import (
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"gorm.io/gorm"
	"math"
	"sort"
	"strings"
)

type MatchmakingController struct {
	userRepo *repositories.UserRepository
	db       *gorm.DB
}

func NewMatchmakingController(userRepo *repositories.UserRepository, db *gorm.DB) *MatchmakingController {
	return &MatchmakingController{userRepo: userRepo, db: db}
}

func (c *MatchmakingController) FindMatches(playerFormId uint, limit int) ([]models.PlayerForm, error) {
	var currentPlayer models.PlayerForm
	if err := c.db.Preload("Statistics").First(&currentPlayer).Error; err != nil {
		return nil, err
	}

	var candidates []models.PlayerForm
	query := c.db.Preload("Statistics").
		Where("player_form_id = ?", playerFormId).
		Where("region = ?", currentPlayer.Region.Region).
		Where("language = ?", currentPlayer.Region.Language).
		Where("array_overlaps(language, ?)", currentPlayer.Region.Language).
		Where("id != ?", playerFormId).
		Limit(limit)
	if err := query.Find(&candidates).Error; err != nil {
		return nil, err
	}
	ranked := c.rankCandidates(currentPlayer, candidates)
	return ranked, nil
}

func (c *MatchmakingController) rankCandidates(currentPlayer models.PlayerForm, candidates []models.PlayerForm) []models.PlayerForm {
	type scoredPlayer struct {
		Player models.PlayerForm
		Score  float64
	}
	var scored []scoredPlayer
	currentStats := currentPlayer.Statistics

	for _, candidate := range candidates {
		score := 0.0

		rankScore := c.compareRanks(currentStats.Rank, candidate.Statistics.Rank)
		score += rankScore * 0.5

		kdScore := 1 - math.Abs(float64(currentStats.KD-candidate.Statistics.KD))/5.0
		score += kdScore * 0.2

		wrScore := 1 - math.Abs(float64(currentStats.WinRate-candidate.Statistics.WinRate))/100.0
		score += wrScore * 0.2

		scored = append(scored, scoredPlayer{candidate, score})
	}
	// сортировка кандидатов
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})
	result := make([]models.PlayerForm, len(scored))
	for i, sp := range scored {
		result[i] = sp.Player
	}
	return result
}

func (c *MatchmakingController) compareRanks(rank1, rank2 string) float64 {
	rankOrder := map[string]int{
		"Iron": 0, "Bronze": 1, "Silver": 2, "Gold": 3,
		"Platinum": 4, "Diamond": 5, "Immortal": 6, "Radiant": 7,
	}
	tier1 := strings.Split(rank1, " ")[0]
	tier2 := strings.Split(rank2, " ")[0]
	return 1.0 - math.Abs(float64(rankOrder[tier1]-rankOrder[tier2]))/7.0
}
