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
	ExchangeRequestStatusPending    ExchangeRequestStatus = "pending"
	ExchangeRequestStatusMatched    ExchangeRequestStatus = "matched"
	ExchangeRequestStatusAccepted   ExchangeRequestStatus = "accepted"
	ExchangeRequestStatusRejected   ExchangeRequestStatus = "recjected"
	ExchangeRequestStatusFoundMatch ExchangeRequestStatus = "found match"
)

func (s ExchangeRequestStatus) String() string {
	return string(s)
}

func (s ExchangeRequestStatus) Badge() string {
	switch s {
	case ExchangeRequestStatusPending:
		return "neutral"
	case ExchangeRequestStatusMatched, ExchangeRequestStatusFoundMatch:
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
	UserEmail     string
	UserGoogleId  string        `gorm:"not null,foreignKey:UserGoogleId,form:user_id"`
	User          User          `gorm:"foreignKey:UserGoogleId;references:GoogleId"`
	DesiredBookID string        `gorm:"not null,foreignKey:BookID" form:"book_id"`
	DesiredBook   Book          `gorm:"foreignKey:DesiredBookID;constraint:OnDelete:SET NULL"`
	OfferedBooks  []OfferedBook `gorm:"constraint:OnDelete:SET NULL"`
	Status        ExchangeRequestStatus
	Matches       []ExchangeMatch `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

func (r *ExchangeRequest) GetMatchStatus(otherRequestId uint) ExchangeRequestStatus {
	for _, match := range r.Matches {
		if match.MatchRequestID == otherRequestId {
			return match.Status
		}
	}
	return ""
}

type ExchangeMatch struct {
	gorm.Model
	ExchangeRequestID uint
	Request           ExchangeRequest `gorm:"foreignKey:ExchangeRequestID"`
	MatchRequestID    uint
	MatchRequest      ExchangeRequest `gorm:"foreignKey:MatchRequestID"`
	Status            ExchangeRequestStatus
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
