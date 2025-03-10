package model

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ID         uint   `json:"id" gorm:"primaryKey"`
	SenderID   uint   `json:"sender_id" gorm:"not null"`
	ReceiverID uint   `json:"receiver_id" gorm:"not null"`
	Content    string `json:"content" gorm:"type:text;not null"`
	Timestamp  int64  `json:"timestamp" gorm:"not null"`
	Read       bool   `json:"read" gorm:"default:false"`
}
type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(100)" json:"name"`
	Email    string `gorm:"uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(100)" json:"password"`
	Online   bool   `json:"online"`
}
