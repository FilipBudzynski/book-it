package models

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	GoogleId string `gorm:"primaryKey" json:"google_id"`
	Books    []UserBook
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
