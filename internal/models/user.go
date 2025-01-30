package models

import (
	"errors"
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

const UserGenresLimit = 5

type User struct {
	GoogleId         string `gorm:"primaryKey;column:google_id" json:"google_id"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt    `gorm:"index"`
	Username         string            `gorm:"not null" json:"username"`
	Email            string            `gorm:"unique;not null;" json:"email"`
	Books            []UserBook        `gorm:"foreignKey:UserGoogleId;constraint:OnDelete:CASCADE;"` // Ensure CASCADE delete on UserBook
	ExchangeRequests []ExchangeRequest `gorm:"foreignKey:UserGoogleId;constraint:OnDelete:CASCADE;"`
	Genres           []Genre           `gorm:"many2many:user_genres;"`
	AvatarURL        string
	Location         *Location `gorm:"foreignKey:UserGoogleId;constraint:OnDelete:CASCADE;"`
}

type Location struct {
	gorm.Model
	UserGoogleId string `gorm:"not null;"`
	Latitude     float64
	Longitude    float64
	Formatted    string
}

var (
	ErrUsernameRequired        = errors.New("username is required")
	ErrEmailRequired           = errors.New("a valid email is required")
	ErrGoogleIdRequired        = errors.New("google id is required")
	ErrUserGenresLimitExceeded = errors.New("user can pick at most 5 genres")
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

func (u *User) HasGenre(genre string) bool {
	genreNames := make([]string, len(u.Genres))
	for i, genre := range u.Genres {
		genreNames[i] = genre.Name
	}
	return slices.Contains(genreNames, genre)
}
