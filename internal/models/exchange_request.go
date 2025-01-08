package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrExchangeRequestNoOfferedBooksProvided = errors.New("no offered books provided in the request")
	ErrExchangeRequestNoDesiredBookProvided  = errors.New("no desired book provided in the request")
)

type ExchangeRequest struct {
	gorm.Model
	UserGoogleId  string `form:"user_id"`
	User          User   `gorm:"foreignKey:UserGoogleId"`
	DesiredBookID string `gorm:"not null,foreignKey:BookID" form:"book_id"`
	DesiredBook   Book   `gorm:"foreignKey:DesiredBookID;constraint:OnDelete:SET NULL"`
	OfferedBooks  []OfferedBook
}

func (e *ExchangeRequest) Validate() error {
	if len(e.OfferedBooks) == 0 {
		return ErrExchangeRequestNoOfferedBooksProvided
	}

	if e.DesiredBookID == "" {
		return ErrExchangeRequestNoDesiredBookProvided
	}

	return nil
}

type OfferedBook struct {
	gorm.Model
	ExchangeRequestID uint
	BookId            string `gorm:"not null,foreignKey:BookID" form:"book_id"`
}
