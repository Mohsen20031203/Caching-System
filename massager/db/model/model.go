package models

import (
	"gorm.io/gorm"
)

type Message struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	SenderNumber   uint           `json:"sender_number" gorm:"not null;index"`
	ReceiverNumber uint           `json:"receiver_number" gorm:"not null;index"`
	Content        string         `json:"content" gorm:"type:text;not null"`
	Read           bool           `json:"read" gorm:"default:false"`
	Status         string         `json:"status" gorm:"type:varchar(20);default:'sent'"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	CreatedAt      int64          `json:"created_at"`
}

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"type:varchar(100);not null"`
	Phone        string `json:"phone" gorm:"type:varchar(20);unique;not null"`
	PasswordHash string `json:"-" gorm:"type:varchar(255);not null"`
	Online       bool   `json:"online" gorm:"default:false"`

	Bio    string `json:"bio" gorm:"type:text"`
	Avatar string `json:"avatar" gorm:"type:varchar(255)"`

	SentMessages     []Message `gorm:"foreignKey:SenderID"`
	ReceivedMessages []Message `gorm:"foreignKey:ReceiverID"`
}
