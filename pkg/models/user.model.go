package models

import (
	"gorm.io/gorm"
)

var MigrateModels = []any{
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

// UserBook model is an abstraction for link between user and a book
// It bounds book to a user providing more information about user interactions with a book
type UserBook struct {
	gorm.Model
	Status       userBookStatus `gorm:"not null"`
	UserGoogleId string         `gorm:"not null"` // foreignKey
	BookId       string         `gorm:"not null"`
}
