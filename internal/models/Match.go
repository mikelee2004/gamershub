package models

import "time"

type Match struct {
	Id     uint      `json:"id" binding:"required"`
	Date   time.Time `json:"timestamptz" binding:"required"`
	UserId uint      `json:"userid" binding:"required"`
	AllyId uint      `json:"allyid" binding:"required"`
	Result uint      `json:"result" binding:"required"`
	Rating int       `json:"rating" binding:"required"`

	// todo
}
