package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username" form:"username"`
	Email    string `gorm:"unique;not null" json:"email" form:"email"`
	ID       uint   `gorm:"primaryKey" json:"id"`
}
