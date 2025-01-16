package models

import (
	"errors"
	"slices"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username         string            `gorm:"not null" json:"username"`
	Email            string            `gorm:"unique;not null" json:"email"`
	GoogleId         string            `gorm:"primaryKey" json:"google_id"`
	Books            []UserBook        `gorm:"onDelete:CASCADE"`
	ExchangeRequests []ExchangeRequest `gorm:"foreignKey:UserGoogleId;onDelete:CASCADE"`
	Genres           []Genre           `gorm:"many2many:user_genres;"`
}

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrEmailRequired    = errors.New("a valid email is required")
	ErrGoogleIdRequired = errors.New("google id is required")
)

func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return ErrUsernameRequired
	}

	if !strings.Contains(u.Email, "@") || len(strings.TrimSpace(u.Email)) < 5 {
		return ErrEmailRequired
	}

	if strings.TrimSpace(u.GoogleId) == "" {
		return ErrGoogleIdRequired
	}

	return nil
}

func (u *User) ContainsGenre(genre string) bool {
	genreNames := make([]string, len(u.Genres))
	for i, genre := range u.Genres {
		genreNames[i] = genre.Name
	}
	return slices.Contains(genreNames, genre)
}
