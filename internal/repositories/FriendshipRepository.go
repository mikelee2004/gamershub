package repositories

import (
	"errors"
	"fmt"
	"gamershub/internal/models"
	"gorm.io/gorm"
)

type FriendshipRepository struct {
	db *gorm.DB
}

func NewFriendshipRepository(db *gorm.DB) *FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) SendRequest(senderID, friendID uint) error {
	if senderID == friendID {
		return errors.New("cannot add yourself as friend")
	}

	// Проверяем, не существует ли уже такой связи
	var existing models.Friendship
	result := r.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		senderID, friendID, friendID, senderID).
		First(&existing)

	if result.Error == nil {
		return errors.New("friendship already exists")
	}

	friendship := &models.Friendship{
		UserId:   senderID,
		FriendId: friendID,
		Status:   models.StatusPending,
	}

	return r.db.Create(friendship).Error
}

// Принять запрос в друзья
func (r *FriendshipRepository) AcceptRequest(userID, friendID uint) error {
	// 1. Находим запрос на дружбу
	var request models.Friendship
	result := r.db.Where("user_id = ? AND friend_id = ? AND status = ?",
		friendID,
		userID,
		models.StatusPending).
		First(&request)

	// 2. Проверяем ошибки
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("friend request not found or already processed")
		}
		return fmt.Errorf("database error: %w", result.Error)
	}

	// 3. Обновляем статус
	if err := r.db.Model(&request).
		Update("status", models.StatusApproved).Error; err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	return nil
}

// Получить список друзей
func (r *FriendshipRepository) GetFriends(userID uint) ([]models.User, error) {
	var friends []models.User

	err := r.db.Joins("JOIN friendships ON users.id = friendships.friend_id").
		Where("friendships.user_id = ? AND friendships.status = ?", userID, models.StatusApproved).
		Or("friendships.friend_id = ? AND friendships.status = ?", userID, models.StatusApproved).
		Find(&friends).Error

	return friends, err
}

// Получить входящие запросы
func (r *FriendshipRepository) GetPendingRequests(userID uint) ([]models.User, error) {
	var users []models.User

	err := r.db.Joins("JOIN friendships ON users.id = friendships.user_id").
		Where("friendships.friend_id = ? AND friendships.status = ?", userID, models.StatusPending).
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("database error3: %w", err)
	}
	return users, nil
}
