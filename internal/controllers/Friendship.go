package controllers

import (
	"gamershub/internal/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type FriendshipController struct {
	friendshipRepo *repositories.FriendshipRepository
	userRepo       *repositories.UserRepository
}

func NewFriendshipController(fr *repositories.FriendshipRepository, ur *repositories.UserRepository) *FriendshipController {
	return &FriendshipController{
		friendshipRepo: fr,
		userRepo:       ur,
	}
}

func (fc *FriendshipController) SendFriendRequest(c *gin.Context) {
	keys := c.Keys
	log.Printf("[DEBUG] Context keys: %v", keys)

	senderID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Приводим senderID к uint
	uid, ok := senderID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Получаем ID получателя из параметров URL
	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID format"})
		return
	}

	// Проверяем, что пользователь не отправляет запрос сам себе
	if uid == uint(friendID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot send friend request to yourself"})
		return
	}

	// Проверяем существование пользователя, которому отправляем запрос
	if _, err := fc.userRepo.FindUserByID(uint(friendID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend user not found"})
		return
	}

	// Отправляем запрос
	if err := fc.friendshipRepo.SendRequest(uid, uint(friendID)); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

func (fc *FriendshipController) AcceptFriendRequest(c *gin.Context) {
	currentUserID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Приводим ID к uint
	uid, ok := currentUserID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
		return
	}

	// Получаем ID отправителя заявки из URL
	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID format"})
		return
	}

	// Принимаем заявку
	if err := fc.friendshipRepo.AcceptRequest(uid, uint(friendID)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted successfully"})
}

func (fc *FriendshipController) GetFriends(c *gin.Context) {
	userID := c.GetUint("userID")

	friends, err := fc.friendshipRepo.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get friends"})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (fc *FriendshipController) GetPendingRequests(c *gin.Context) {
	userID := c.GetUint("userID")

	requests, err := fc.friendshipRepo.GetPendingRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get requests"})
		return
	}

	c.JSON(http.StatusOK, requests)
}
