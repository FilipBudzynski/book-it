package models

import (
	"errors"

	"gorm.io/gorm"
)

type (
	Status interface {
		String() string
		Badge() string
	}

	ExchangeRequestStatus string
)

const (
	ExchangeRequestStatusPending    ExchangeRequestStatus = "pending"
	ExchangeRequestStatusMatched    ExchangeRequestStatus = "matched"
	ExchangeRequestStatusAccepted   ExchangeRequestStatus = "accepted"
	ExchangeRequestStatusRejected   ExchangeRequestStatus = "recjected"
	ExchangeRequestStatusFoundMatch ExchangeRequestStatus = "found match"
	ExchangeRequestStatusActive     ExchangeRequestStatus = "active"
)

func (s ExchangeRequestStatus) String() string {
	return string(s)
}

func (s ExchangeRequestStatus) Badge() string {
	switch s {
	case ExchangeRequestStatusPending:
		return "neutral"
	case ExchangeRequestStatusMatched,
		ExchangeRequestStatusFoundMatch:
		return "info"
	case ExchangeRequestStatusAccepted,
		ExchangeRequestStatusActive:
		return "success"
	case ExchangeRequestStatusRejected:
		return "error"
	}
	return "secondary"
}

var (
	ErrExchangeRequestNoOfferedBooksProvided      = errors.New("no offered books provided in the request")
	ErrExchangeRequestNoDesiredBookProvided       = errors.New("no desired book provided in the request")
	ErrExchangeRequestDuplicateOfferedBooks       = errors.New("duplicate offered books in the request")
	ErrExchangeRequestActiveRequestWithThisBookID = errors.New("exchange request with this book id is active")
	ErrExchangeRequestLatitudeOutOfRange          = errors.New("latitude out of range")
	ErrExchangeRequestLongitudeOutOfRange         = errors.New("longitude out of range")
)

type ExchangeRequest struct {
	gorm.Model
	UserEmail     string
	UserGoogleId  string        `gorm:"type:text;not null"`
	User          User          
	DesiredBookID string        `gorm:"not null,foreignKey:BookID" form:"book_id"`
	DesiredBook   Book          `gorm:"foreignKey:DesiredBookID;constraint:OnDelete:SET NULL"`
	OfferedBooks  []OfferedBook `gorm:"constraint:OnDelete:SET NULL"`
	Status        ExchangeRequestStatus
	Matches       []ExchangeMatch `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Latitude      float64
	Longitude     float64
}

func (r *ExchangeRequest) GetMatchStatus(otherRequestId uint) MatchStatus {
	for _, match := range r.Matches {
		if match.MatchedExchangeRequestID == otherRequestId || match.ExchangeRequestID == otherRequestId {
			return match.Status
		}
	}
	return MatchStatusPending
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

	if e.Latitude < -90 || e.Latitude > 90 {
		return ErrExchangeRequestLatitudeOutOfRange
	}
	if e.Longitude < -180 || e.Longitude > 180 {
		return ErrExchangeRequestLongitudeOutOfRange
	}

	return nil
}

func (e *ExchangeRequest) checkDuplicates() error {
	seenBooks := make(map[string]bool)
	for _, book := range e.OfferedBooks {
		if seenBooks[book.BookID] {
			return ErrExchangeRequestDuplicateOfferedBooks
		}
		seenBooks[book.BookID] = true
	}
	return nil
}

type OfferedBook struct {
	gorm.Model
	ExchangeRequestID uint
	BookID            string `gorm:"not null,foreignKey:BookID" form:"book_id"`
	Book              Book   `gorm:"foreignKey:BookID;constraint:OnDelete:SET NULL"`
}
