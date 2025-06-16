package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gamershub/internal/models"
	"gamershub/internal/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
)

// MatchmakingController handles matchmaking operations
type MatchmakingController struct {
	playerForm *repositories.FormRepository
	db         *gorm.DB
}

func NewMatchmakingController(playerForm *repositories.FormRepository, db *gorm.DB) *MatchmakingController {
	return &MatchmakingController{playerForm: playerForm, db: db}
}

// GetMatches godoc
// @Summary Find matching players
// @Description Finds potential gaming partners based on player preferences and statistics
// @Tags Matchmaking
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} models.MatchResponse "List of matching players"
// @Success 200 {object} models.NoMatchesResponse "When no matches found"
// @Failure 404 {object} models.ErrorResponse "Player form not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /lft/find [get]
func (m *MatchmakingController) GetMatches(c *gin.Context) {
	currentUserID := c.GetUint("userID")

	var currentPlayer models.PlayerForm
	if err := m.db.Where("user_id = ?", currentUserID).First(&currentPlayer).Error; err != nil {
		log.Printf("Error loading current player: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Player form not found"})
		return
	}
	log.Printf("Loaded current player from DB: %+v", currentPlayer)

	var candidates []models.PlayerForm
	if err := m.db.Where("user_id != ?", currentUserID).Find(&candidates).Error; err != nil {
		log.Printf("Error loading candidates: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not load players"})
		return
	}
	log.Printf("Loaded %d candidates from DB", len(candidates))

	for i, cand := range candidates {
		log.Printf("Candidate %d: %+v", i+1, cand)
	}
	criteria := models.MatchCriteria{
		PlayerForm:  currentPlayer,
		MaxDistance: 300,
	}

	matches, err := m.FindMatches(criteria, candidates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Matching failed",
			"details": err,
		})
		return
	}

	if matches == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No matches found",
			"details": "Try adjusting your search criteria",
			"your_criteria": gin.H{
				"region":    currentPlayer.Region.Region,
				"languages": currentPlayer.Region.Language,
				"rank":      currentPlayer.Statistics.Rank,
				"goal":      currentPlayer.Goal,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count":   matches,
		"matches": matches,
	})
}

// SendInvite godoc
// @Summary Send game invitation
// @Description Send a game invitation to another player
// @Tags Matchmaking
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body models.SendInviteRequest true "Invitation details"
// @Success 201 {object} models.InvitationSentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /lft/invite [post]
func (m *MatchmakingController) SendInvite(c *gin.Context) {
	currentUserID := c.GetUint("userID")
	var req struct {
		RecipientID uint   `json:"recipient_id" binding:"required"`
		Message     string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var recipient models.PlayerForm
	if err := m.db.Where("user_id = ?", req.RecipientID).First(&recipient).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipient not found"})
		return
	}
	var existingInvitation models.GameSession
	if err := m.db.Where("sender_id = ? AND recipient_id = ? AND status = 'pending'",
		currentUserID, req.RecipientID).First(&existingInvitation).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Invitation already sent"})
		return
	}

	invitation := models.GameSession{
		SenderID:    currentUserID,
		RecipientID: req.RecipientID,
		Message:     req.Message,
		Status:      "pending",
		ExpiresAt:   time.Now().Add(24 * time.Hour), // Приглашение активно 24 часа
	}
	if err := m.db.Create(&invitation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send invitation"})
		return
	}
	//	TODO: WebSocket Push-Уведомления/Email-оповещение

	c.JSON(http.StatusCreated, gin.H{
		"message":       "Invitation sent successfully",
		"invitation_id": invitation.ID,
	})
}

// GetInvites godoc
// @Summary Get pending invitations
// @Description Get list of pending game invitations
// @Tags Matchmaking
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.InvitationResponse "List of invitations"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /lft/invites [get]
func (m *MatchmakingController) GetInvites(c *gin.Context) {
	userID := c.GetUint("userID")

	var invitations []models.GameSession
	if err := m.db.Where("recipient_id = ? AND status = 'pending' AND expires_at > ?",
		userID, time.Now()).Find(&invitations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load invitations"})
		return
	}

	var result []map[string]interface{}
	for _, inv := range invitations {
		var sender models.PlayerForm
		if err := m.db.Where("user_id = ?", inv.SenderID).First(&sender).Error; err != nil {
			continue
		}

		// Десериализуем Statistics из JSON
		var stats models.Statistics
		if err := json.Unmarshal([]byte(sender.StatisticsJSON), &stats); err != nil {
			log.Printf("Failed to unmarshal stats for user %d: %v", sender.UserID, err)
			continue
		}

		result = append(result, map[string]interface{}{
			"invitation_id": inv.ID,
			"sender": map[string]interface{}{
				"user_id":     sender.UserID,
				"valorant_id": sender.ValorantID,
				"rank":        stats.Rank,
				"rank_name":   m.getRankName(stats.Rank),
			},
			"message":    inv.Message,
			"sent_at":    inv.CreatedAt,
			"expires_at": inv.ExpiresAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"invitations": result})
}

// RespondToInvitation godoc
// @Summary Respond to game invitation
// @Description Accept or decline game invitation
// @Tags Matchmaking
// @Accept json
// @Produce json
// @Param request body models.RespondInvitationRequest true "Response data"
// @Success 200 {object} models.InvitationResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /lft/respond [post]
func (m *MatchmakingController) RespondToInvitation(c *gin.Context) {
	userID := c.GetUint("userID")
	var req struct {
		InvitationID uint   `json:"invitation_id" binding:"required"`
		Action       string `json:"action" binding:"required,oneof=accept decline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}
	// Начинаем транзакцию
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Находим приглашение с блокировкой
	var invitation models.GameSession
	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("id = ? AND recipient_id = ?", req.InvitationID, userID).
		First(&invitation).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Invitation not found or already processed"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if invitation.Status != "pending" {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{
			"error":          "Invitation already processed",
			"current_status": invitation.Status,
		})
		return
	}

	if time.Now().After(invitation.ExpiresAt) {
		if err := tx.Model(&invitation).Update("status", "expired").Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update invitation status"})
			return
		}
		tx.Commit()
		c.JSON(http.StatusGone, gin.H{"error": "Invitation expired"})
		return
	}

	switch req.Action {
	case "accept":
		var sender models.User
		if err := tx.First(&sender, invitation.SenderID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusConflict, gin.H{"error": "Invitation sender not found"})
			return
		}
		updates := map[string]interface{}{
			"status":      "accepted",
			"accepted_at": time.Now(),
			"started_at":  time.Now(),
			"updated_at":  time.Now(),
		}
		if err := tx.Model(&invitation).Updates(updates).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to update invitation: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to accept invitation",
				"details": err.Error(),
			})
			return
		}

		if err := tx.Model(&models.GameSession{}).
			Where("((sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)) AND status = 'pending'",
				invitation.SenderID, userID, userID, invitation.SenderID).
			Updates(map[string]interface{}{
				"status":     "canceled",
				"updated_at": time.Now(),
			}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel other invitations"})
			return
		}
	case "decline":
		if err := tx.Model(&invitation).Updates(map[string]interface{}{
			"status":     "declined",
			"updated_at": time.Now(),
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline invitation"})
			return
		}
	}
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    fmt.Sprintf("Invitation %sd", req.Action),
		"invitation": invitation,
	})
}

func (m *MatchmakingController) FindMatches(criteria models.MatchCriteria, allPlayers []models.PlayerForm) (interface{}, interface{}) {
	if len(allPlayers) == 0 {
		return nil, errors.New("no players to match with")
	}
	if criteria.MaxDistance == 0 {
		criteria.MaxDistance = 300
	}

	var results []models.MatchResult

	if err := criteria.PlayerForm.DeserializeJSONData(); err != nil {
		log.Printf("Failed to deserialize player data for user %d: %v", criteria.PlayerForm.UserID, err)
		return nil, err
	}

	for i, candidate := range allPlayers {
		if candidate.UserID == criteria.PlayerForm.UserID {
			continue
		}
		if err := candidate.DeserializeJSONData(); err != nil {
			log.Printf("Failed to deserialize candidate %d data: %v", candidate.UserID, err)
			continue
		}
		log.Printf("Checking candidate %d/%d: UserID %d (Rank: %d, Region: %s, Goal: %s)",
			i+1, len(allPlayers),
			candidate.UserID,
			candidate.Statistics.Rank,
			candidate.Region.Region,
			candidate.Goal)
		// проверка базовых критериев
		passed, reason := checkBasicCriteriaWithReason(criteria.PlayerForm, candidate, criteria.MaxDistance)
		if !passed {
			log.Printf("Candidate %d rejected: %s", candidate.UserID, reason)
			continue
		}
		// оценка совпадения
		score := calculateMatchScore(criteria.PlayerForm, candidate)
		log.Printf("Candidate %d match score: %.2f", candidate.UserID, score)
		if score > 0 {
			results = append(results, models.MatchResult{
				Player:     candidate,
				MatchScore: score,
			})
		}
	}

	// сортировка
	sort.Slice(results, func(i, j int) bool {
		return results[i].MatchScore > results[j].MatchScore
	})
	log.Printf("Found %d matches for user %d", len(results), criteria.PlayerForm.UserID)
	return results, nil
}

// checkBasicCriteria проверяет обязательные критерии
func checkBasicCriteriaWithReason(player, candidate models.PlayerForm, maxDistance int) (bool, string) {
	if player.Statistics.Rank == 0 || candidate.Statistics.Rank == 0 {
		return false, "rank data not loaded"
	}
	if player.Region.Region == "" || candidate.Region.Region == "" {
		return false, "region data not loaded"
	}
	if !strings.EqualFold(player.Region.Region, candidate.Region.Region) {
		return false, fmt.Sprintf("region mismatch (%s != %s)",
			player.Region.Region, candidate.Region.Region)
	}
	if len(player.Region.Language) == 0 || len(candidate.Region.Language) == 0 {
		return false, "no languages specified"
	}
	hasCommonLanguage := false
	for _, plang := range player.Region.Language {
		for _, clang := range candidate.Region.Language {
			if strings.EqualFold(plang, clang) {
				hasCommonLanguage = true
				break
			}
		}
		if hasCommonLanguage {
			break
		}
	}
	if !hasCommonLanguage {
		return false, fmt.Sprintf("no common language (%v vs %v)",
			player.Region.Language, candidate.Region.Language)
	}
	rankDiff := abs(int(player.Statistics.Rank) - int(candidate.Statistics.Rank))
	if rankDiff > maxDistance {
		return false, fmt.Sprintf("rank difference too big (%d > %d)",
			rankDiff, maxDistance)
	}
	return true, ""
}

// calculateMatchScore рассчитывает оценку совпадения (0-1)
func calculateMatchScore(player, candidate models.PlayerForm) float64 {
	var score float64
	maxScore := 0.0

	// вес 0.3
	if player.Goal == candidate.Goal {
		score += 0.3
	}
	maxScore += 0.3

	// вес 0.25
	rankDiff := abs(int(player.Statistics.Rank) - int(candidate.Statistics.Rank))
	rankScore := 1.0 - float64(rankDiff)/300.0
	if rankScore < 0 {
		rankScore = 0
	}
	score += rankScore * 0.25
	maxScore += 0.25

	// вес 0.2
	commonComms := 0
	for _, way := range player.CommunicationWays {
		for _, candidateWay := range candidate.CommunicationWays {
			if way == candidateWay {
				commonComms++
				break
			}
		}
	}
	if len(player.CommunicationWays) > 0 {
		commScore := float64(commonComms) / float64(len(player.CommunicationWays))
		score += commScore * 0.2
	}
	maxScore += 0.2

	// Схожесть статистики (15% веса)
	statsScore := float64(calculateStatsSimilarity(player.Statistics, candidate.Statistics))
	score += statsScore * 0.15
	maxScore += 0.15

	// Описание (10% веса) - можно добавить анализ текста, но пока просто наличие
	if player.Description != "" && candidate.Description != "" {
		score += 0.1
	}
	maxScore += 0.1

	// нормализация оценки
	if maxScore > 0 {
		return score / maxScore
	}
	return 0
}

// calculateStatsSimilarity рассчитывает схожесть статистики
func calculateStatsSimilarity(player, candidate models.Statistics) float32 {
	var similarity float32

	// KD ratio (30% веса)
	kdDiff := absFloat(player.KD - candidate.KD)
	kdScore := 1.0 - kdDiff/2.0
	if kdScore < 0 {
		kdScore = 0
	}
	similarity += kdScore * 0.3

	// Win rate (30% веса)
	wrDiff := absFloat(player.WinRate - candidate.WinRate)
	wrScore := 1.0 - wrDiff/100.0
	if wrScore < 0 {
		wrScore = 0
	}
	similarity += wrScore * 0.3

	// HS percent (20% веса)
	hsDiff := absFloat(player.HSPercent - candidate.HSPercent)
	hsScore := 1.0 - hsDiff/100.0
	if hsScore < 0 {
		hsScore = 0
	}
	similarity += hsScore * 0.2

	// Peak rank (20% веса)
	peakDiff := abs(int(player.PeakRank) - int(candidate.PeakRank))
	peakScore := 1.0 - float32(peakDiff)/300.0
	if peakScore < 0 {
		peakScore = 0
	}
	similarity += peakScore * 0.2

	return similarity
}

func (m *MatchmakingController) getRankName(rankID uint) string {
	var rank models.Rank
	if err := m.db.Where("id = ?", rankID).First(&rank).Error; err != nil {
		return "Unknown"
	}
	return rank.Name
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func absFloat(n float32) float32 {
	if n < 0 {
		return -n
	}
	return n
}
