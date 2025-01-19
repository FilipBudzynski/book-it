package models

import "gorm.io/gorm"

type (
	MatchStatus   string
	MatchDecision string
)

const (
	MatchStatusPending  MatchStatus = "pending"
	MatchStatusAccepted MatchStatus = "accepted"
	MatchStatusDeclined MatchStatus = "declined"

	MatchDecisionPending  MatchDecision = "pending"
	MatchDecisionAccepted MatchDecision = "accepted"
	MatchDecisionDeclined MatchDecision = "declined"
)

func (s MatchDecision) Accepted() bool {
	return s == MatchDecisionAccepted
}

func (s MatchDecision) Declined() bool {
	return s == MatchDecisionDeclined
}

func (s MatchStatus) String() string {
	return string(s)
}

func (s MatchStatus) Badge() string {
	switch s {
	case MatchStatusPending:
		return "info"
	case MatchStatusAccepted:
		return "success"
	case MatchStatusDeclined:
		return "error"
	}
	return "secondary"
}

type ExchangeMatch struct {
	gorm.Model
	ExchangeRequestID        uint            `gorm:"uniqueIndex:id_match"`
	Request                  ExchangeRequest `gorm:"foreignKey:ExchangeRequestID"`
	MatchedExchangeRequestID uint            `gorm:"uniqueIndex:id_match"`
	MatchedExchangeRequest   ExchangeRequest `gorm:"foreignKey:MatchedExchangeRequestID"`
	Request1Decision         MatchDecision
	Request2Decision         MatchDecision
	Status                   MatchStatus
    Distance                 float64
}

func (e *ExchangeMatch) MatchedRequest(requestId uint) *ExchangeRequest {
	if e.ExchangeRequestID == requestId {
		return &e.MatchedExchangeRequest
	}
	return &e.Request
}

func (e *ExchangeMatch) IsAccepted(requestId uint) bool{
    if e.ExchangeRequestID == requestId {
        return e.Request1Decision.Accepted()
    } else {
        return e.Request2Decision.Accepted()
    }
}

func (e *ExchangeMatch) GetDecision(requestId uint) MatchDecision {
    if e.ExchangeRequestID == requestId {
        return e.Request1Decision
    }
    return e.Request2Decision
}

