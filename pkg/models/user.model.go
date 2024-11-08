package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"size:255" json:"password"`
	GoogleId string `gorm:"size:255" json:"google_id"`
	// CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP()" json:"created_at"`
	ID uint `gorm:"primaryKey" json:"id"`
}
