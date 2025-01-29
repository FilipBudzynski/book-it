package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrUserBookQueryWithoutId          = errors.New("user book ID not provided in query parameters")
	ErrUserBookInActiveExchangeRequest = errors.New("user book in active exchange request")
)

// UserBook model is an abstraction for link between user and a book
// It bounds book to a user providing more information about user interactions with a book
type UserBook struct {
	gorm.Model
	UserGoogleId    string           `gorm:"not null;"` // foreignKey
	BookID          string           `gorm:"not null;"`
	Book            Book             `gorm:"foreignKey:BookID;constraint:OnDelete:SET NULL;"`
	ReadingProgress *ReadingProgress `gorm:"foreignKey:UserBookID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func BookInUserBooks(bookID string, userBooks []*UserBook) bool {
	for _, userBook := range userBooks {
		if userBook.BookID == bookID {
			return true
		}
	}
	return false
}
