package models

import (
	"gorm.io/gorm"
)

var Models = []any{
	&User{},
	&UserBook{},
}

type userBookStatus string

const (
	BookStatusInProgress userBookStatus = "InProgress"
	BookStatusNotStarted userBookStatus = "NotStarted"
	BookStatusStarted    userBookStatus = "Started"
	BookStatusFinished   userBookStatus = "Finished"
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
	ID            string
	ISBN          uint     `gorm:"primaryKey"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Description   string   `json:"description"`
	ImageLink     string   `json:"thumbnail"`
	Genre         string
	Link          string
	PublishedDate string
}

type UserBook struct {
	gorm.Model
	Status       userBookStatus `gorm:"not null"`
	UserGoogleId string         `gorm:"not null"` // foreignKey
	BookId       string         `gorm:"not null"`
}
