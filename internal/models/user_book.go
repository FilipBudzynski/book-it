package models

import (
	"errors"

	"gorm.io/gorm"
)

var ErrUserBookQueryWithoutId = errors.New("something went wrong with the request. Book ID was not provided in query parameters")

// UserBook model is an abstraction for link between user and a book
// It bounds book to a user providing more information about user interactions with a book
type UserBook struct {
	gorm.Model
	UserGoogleId    string `gorm:"not null"` // foreignKey
	BookID          string `gorm:"not null"`
	Book            Book
	ReadingProgress *ReadingProgress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
