package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	GoogleId string `gorm:"primaryKey" json:"google_id"`
	Books []UserBook
}

type UserBook struct {
	gorm.Model

	Name         string `gorm:"not null`
	UserGoogleId string `gorm:not null`
	ID           uint   `gorm:primaryKet`
}
