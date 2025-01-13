package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrExchangeRequestNoOfferedBooksProvided = errors.New("no offered books provided in the request")
	ErrExchangeRequestNoDesiredBookProvided  = errors.New("no desired book provided in the request")
	ErrExchangeRequestDuplicateOfferedBooks  = errors.New("duplicate offered books in the request")
)

type ExchangeRequestStatus string

const (
	ExchangeRequestStatusPending  ExchangeRequestStatus = "pending"
	ExchangeRequestStatusMatched  ExchangeRequestStatus = "matched"
	ExchangeRequestStatusAccepted ExchangeRequestStatus = "accepted"
	ExchangeRequestStatusRejected ExchangeRequestStatus = "recjected"
)

func (s ExchangeRequestStatus) String() string {
	return string(s)
}

func (s ExchangeRequestStatus) Badge() string {
	switch s {
	case ExchangeRequestStatusPending:
		return "neutral"
	case ExchangeRequestStatusMatched:
		return "info"
	case ExchangeRequestStatusAccepted:
		return "success"
	case ExchangeRequestStatusRejected:
		return "error"
	}
	return "secondary"
}

type ExchangeRequest struct {
	gorm.Model
	UserGoogleId  string        `form:"user_id"`
	User          User          `gorm:"foreignKey:UserGoogleId"`
	DesiredBookID string        `gorm:"not null,foreignKey:BookID" form:"book_id"`
	DesiredBook   Book          `gorm:"foreignKey:DesiredBookID;constraint:OnDelete:SET NULL"`
	OfferedBooks  []OfferedBook `gorm:"constraint:OnDelete:SET NULL"`
	Status        ExchangeRequestStatus
	Matches       []ExchangeMatch `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

type ExchangeMatch struct {
	gorm.Model
	ExchangeRequestID uint
	Request           ExchangeRequest `gorm:"foreignKey:ExchangeRequestID"`
	MatchRequestID    uint
	MatchRequest      ExchangeRequest `gorm:"foreignKey:MatchRequestID"`
}

func (m *ExchangeMatch) AfterSave(db *gorm.DB) error {
	var matchCount int64
	err := db.Model(&ExchangeMatch{}).
		Where("exchange_request_id = ?", m.ExchangeRequestID).
		Count(&matchCount).Error
	if err != nil {
		return err
	}

	if matchCount == 0 {
		err = db.Model(&ExchangeRequest{}).
			Where("id = ?", m.ExchangeRequestID).
			Update("Status", ExchangeRequestStatusPending).Error
		if err != nil {
			return err
		}
	} else {

		err := db.Model(&ExchangeRequest{}).
			Where("id = ?", m.ExchangeRequestID).
			Update("Status", ExchangeRequestStatusMatched).Error
		if err != nil {
			return err
		}

		err = db.Model(&ExchangeRequest{}).
			Where("id = ?", m.MatchRequestID).
			Update("Status", ExchangeRequestStatusMatched).Error
		if err != nil {
			return err
		}
	}

	// Check for the reciprocal match and update if necessary
	var reciprocalMatch ExchangeMatch
	err = db.Where("exchange_request_id = ? AND match_request_id = ?", m.MatchRequestID, m.ExchangeRequestID).
		First(&reciprocalMatch).Error

	// Create the reciprocal match if it doesn't exist
	if err != nil && err == gorm.ErrRecordNotFound {
		reciprocalMatch = ExchangeMatch{
			ExchangeRequestID: m.MatchRequestID,
			MatchRequestID:    m.ExchangeRequestID,
		}
		err = db.Create(&reciprocalMatch).Error
		if err != nil {
			return err
		}
	}

	return nil
}

type OfferedBook struct {
	gorm.Model
	ExchangeRequestID uint
	BookId            string `gorm:"not null,foreignKey:BookID" form:"book_id"`
	Book              Book   `gorm:"foreignKey:BookId;constraint:OnDelete:SET NULL"`
}

func (e *ExchangeRequest) Validate() error {
	if len(e.OfferedBooks) == 0 {
		return ErrExchangeRequestNoOfferedBooksProvided
	}

	if e.DesiredBookID == "" {
		return ErrExchangeRequestNoDesiredBookProvided
	}

	if err := e.checkDuplicates(); err != nil {
		return err
	}

	return nil
}

func (e *ExchangeRequest) checkDuplicates() error {
	seenBooks := make(map[string]bool)
	for _, book := range e.OfferedBooks {
		if seenBooks[book.BookId] {
			return ErrExchangeRequestDuplicateOfferedBooks
		}
		seenBooks[book.BookId] = true
	}
	return nil
}
