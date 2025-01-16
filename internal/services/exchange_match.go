package services

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
)

func (s *exchangeService) CreateMatch(requestId, otherRequestId uint) (*models.ExchangeMatch, error) {
	match := &models.ExchangeMatch{
		ExchangeRequestID:        requestId,
		MatchedExchangeRequestID: otherRequestId,
		Status:                   models.MatchStatusPending,
		Request1Decision:         models.MatchDecisionPending,
		Request2Decision:         models.MatchDecisionPending,
	}
	return match, s.repo.CreateMatch(match)
}

func (s *exchangeService) GetMatches(requestId string) ([]*models.ExchangeMatch, error) {
	return s.repo.GetAllMatches(requestId)
}

func (s *exchangeService) GetMatch(id string) (*models.ExchangeMatch, error) {
	return s.repo.GetMatchByID(id)
}

func (s *exchangeService) AcceptMatch(matchID, requestID string) (*models.ExchangeMatch, error) {
	match, err := s.repo.GetMatch(matchID, requestID)
	if err != nil {
		return nil, err
	}

	if err := updateMatchDecision(match, requestID, models.MatchDecisionAccepted); err != nil {
		return nil, err
	}

	if match.Request1Decision.Accepted() && match.Request2Decision.Accepted() {
		match.Status = models.MatchStatusAccepted
	} else if match.Request1Decision.Declined() || match.Request2Decision.Declined() {
		match.Status = models.MatchStatusDeclined
	}

	if err := s.repo.UpdateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *exchangeService) DeclineMatch(matchID, requestID string) (*models.ExchangeMatch, error) {
	match, err := s.repo.GetMatch(matchID, requestID)
	if err != nil {
		return nil, err
	}

	if err := updateMatchDecision(match, requestID, models.MatchDecisionDeclined); err != nil {
		return nil, err
	}
	match.Status = models.MatchStatusDeclined

	if err := s.repo.UpdateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func updateMatchDecision(match *models.ExchangeMatch, requestID string, decision models.MatchDecision) error {
	requestIDParsed, err := utils.ParseStringToUint(requestID)
	if err != nil {
		return err
	}

	switch requestIDParsed {
	case match.ExchangeRequestID:
		match.Request1Decision = decision
	case match.MatchedExchangeRequestID:
		match.Request2Decision = decision
	}

	return nil
}
