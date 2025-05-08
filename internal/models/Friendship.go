package models

import (
	"gorm.io/gorm"
)

type FriendStatus string

const (
	StatusPending  FriendStatus = "PENDING"
	StatusApproved FriendStatus = "APPROVED"
	StatusRejected FriendStatus = "REJECTED"
	StatusBlocked  FriendStatus = "BLOCKED"
)

type Friendship struct {
	gorm.Model
	UserId   uint         `gorm:"not null"`
	FriendId uint         `gorm:"not null"`
	Status   FriendStatus `gorm:"type:varchar(20);default:'PENDING'"`
}
