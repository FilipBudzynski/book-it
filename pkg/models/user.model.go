package models

import (
	"time"

	"gorm.io/gorm"
)

type userBookStatus string

const (
	InProgress userBookStatus = "InProgress"
	NotStarted userBookStatus = "NotStarted"
	Started    userBookStatus = "Started"
	Finished   userBookStatus = "Finished"
)

type User struct {
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	GoogleId string `gorm:"primaryKey" json:"google_id"`
	Books    []UserBook
	gorm.Model
}

// will this be needed?
// maybe just get the books from external API
type Book struct {
	ISBN          uint     `gorm:"primaryKey"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Description   string   `json:"description"`
	ImageLink     string   `json:"thumbnail"`
	Genre         string
	Link          string
	PublishedDate string
	gorm.Model
}

type UserBook struct {
	Status       userBookStatus `gorm:"not null"`
	UserGoogleId string         `gorm:"not null"` // foreignKey
	BookID       uint           `gorm:"not null"`
	AddedAt      time.Time
	gorm.Model
}
