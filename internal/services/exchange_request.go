package services

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/internal/models"
)

type ExchangeRequestRepository interface {
	Create(exchange *models.ExchangeRequest) error
	Get(id, userId string) (*models.ExchangeRequest, error)
	GetAll(userId string) ([]*models.ExchangeRequest, error)
	Delete(id string) error
	DeleteMatchesForRequest(requestId string) error
	Update(exchange *models.ExchangeRequest) error
	FindMatchingRequests(userId, requestId, desiredBookId string, offeredBooks []string) ([]*models.ExchangeRequest, error)
	// match
	CreateMatch(match *models.ExchangeMatch) error
	GetMatch(requestId, otherRequestId uint) (*models.ExchangeMatch, error)
	UpdateMatch(match *models.ExchangeMatch) error
	GetAllMatches(requestId string) ([]*models.ExchangeMatch, error)
	GetMatchByID(id string) (*models.ExchangeMatch, error)
	// utility
	GetActiveExchangeRequestsByBookID(id string) ([]*models.ExchangeRequest, error)
}

type exchangeService struct {
	repo ExchangeRequestRepository
}

func NewExchangeService(r ExchangeRequestRepository) *exchangeService {
	return &exchangeService{
		repo: r,
	}
}

func (s *exchangeService) Create(userId, userEmail, desiredBookID string, userBookIDs []string) (*models.ExchangeRequest, error) {
	offeredBooks := make([]models.OfferedBook, len(userBookIDs))
	for i, id := range userBookIDs {
		offeredBooks[i] = models.OfferedBook{BookId: id}
	}

	exchange := &models.ExchangeRequest{
		UserGoogleId:  userId,
		UserEmail:     userEmail,
		DesiredBookID: desiredBookID,
		OfferedBooks:  offeredBooks,
		Status:        models.ExchangeRequestStatusActive,
	}

	if err := exchange.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.Create(exchange); err != nil {
		return nil, err
	}
	return s.repo.Get(fmt.Sprintf("%d", exchange.ID), userId)
}

func (s *exchangeService) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	return s.repo.GetAll(userId)
}

func (s *exchangeService) Get(id, userId string) (*models.ExchangeRequest, error) {
	return s.repo.Get(id, userId)
}

func (s *exchangeService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return s.repo.DeleteMatchesForRequest(id)
}

func (s *exchangeService) FindMatchingRequests(requestId, userId string) ([]*models.ExchangeRequest, error) {
	r, err := s.Get(requestId, userId)
	if err != nil {
		return nil, err
	}

	matchingRequests, err := s.repo.FindMatchingRequests(
		requestId,
		userId,
		r.DesiredBookID,
		getOfferedBookIDs(r.OfferedBooks),
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(r); err != nil {
		return nil, err
	}

	// var matches []*models.ExchangeMatch
	for _, matchingReq := range matchingRequests {
		_, err := s.CreateMatch(r.ID, matchingReq.ID)
		if err != nil {
			return nil, err
		}
		// matches = append(matches, match)
	}

	return matchingRequests, nil
}

func (s *exchangeService) GetMatches(requestId string) ([]*models.ExchangeMatch, error) {
	return s.repo.GetAllMatches(requestId)
}

func (s *exchangeService) GetMatch(id string) (*models.ExchangeMatch, error) {
	return s.repo.GetMatchByID(id)
}

func (s *exchangeService) AcceptMatch(matchID, requestId uint) (*models.ExchangeMatch, error) {
	match, err := s.repo.GetMatch(matchID, requestId)
	if err != nil {
		return nil, err
	}

	if match.ExchangeRequestID == requestId {
		match.Request1Decision = models.MatchDecisionAccepted
	} else if match.MatchedExchangeRequestID == requestId {
		match.Request2Decision = models.MatchDecisionAccepted
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

// Helper function to extract BookIDs from OfferedBooks
func getOfferedBookIDs(offeredBooks []models.OfferedBook) []string {
	ids := make([]string, len(offeredBooks))
	for i, book := range offeredBooks {
		ids[i] = book.BookId
	}
	return ids
}

func (s *exchangeService) DeclineMatch(matchID, requestID uint) (*models.ExchangeMatch, error) {
	match, err := s.repo.GetMatch(matchID, requestID)
	if err != nil {
		return nil, err
	}

	if match.ExchangeRequestID == requestID {
		match.Request1Decision = models.MatchDecisionDeclined
	} else if match.MatchedExchangeRequestID == requestID {
		match.Request2Decision = models.MatchDecisionDeclined
	}
	match.Status = models.MatchStatusDeclined

	if err := s.repo.UpdateMatch(match); err != nil {
		return nil, err
	}
	return match, nil
}

func (s *exchangeService) CreateMatch(requestId, otherRequestId uint) (*models.ExchangeMatch, error) {
	match := &models.ExchangeMatch{
		ExchangeRequestID:        requestId,
		MatchedExchangeRequestID: otherRequestId,
		Status:                   models.MatchStatusPending,
	}
	return match, s.repo.CreateMatch(match)
}

func (s *exchangeService) CheckMatch(userReqId, otherReqId uint) (bool, error) {
	match, err := s.repo.GetMatch(userReqId, otherReqId)
	if err != nil {
		return false, err
	}

	var found bool
	_, err = s.repo.GetMatch(otherReqId, userReqId)
	if err != nil {
		match.Status = models.MatchStatusPending
		found = false
	} else {
		match.Status = models.MatchStatusAccepted
		found = true
	}

	if err = s.repo.UpdateMatch(match); err != nil {
		return false, err
	}
	return found, nil
}
