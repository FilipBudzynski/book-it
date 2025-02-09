package services

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/internal/geo"
	"github.com/FilipBudzynski/book_it/internal/models"
)

type ExchangeRequestRepository interface {
	Create(exchange *models.ExchangeRequest) error
	Get(id, userId string) (*models.ExchangeRequest, error)
	GetByID(id string) (*models.ExchangeRequest, error)
	GetAll(userId string) ([]*models.ExchangeRequest, error)
	GetAllWithStatus(userId string, status models.ExchangeRequestStatus) ([]*models.ExchangeRequest, error)
	Delete(id string) error
	DeleteMatchesForRequest(requestId string) error
	Update(exchange *models.ExchangeRequest) error
	FindMatchingRequests(userId, requestId, desiredBookId string, offeredBooks []string) ([]*models.ExchangeRequest, error)
	CreateMatch(match *models.ExchangeMatch) error
	GetMatch(requestId, otherRequestId string) (*models.ExchangeMatch, error)
	UpdateMatch(match *models.ExchangeMatch) error
	GetAllMatches(requestId string) ([]*models.ExchangeMatch, error)
	GetMatchByID(id string) (*models.ExchangeMatch, error)
	GetActiveExchangeRequestsByBookID(id string, userID string) ([]*models.ExchangeRequest, error)
}

type exchangeService struct {
	repo ExchangeRequestRepository
}

func NewExchangeService(r ExchangeRequestRepository) *exchangeService {
	return &exchangeService{
		repo: r,
	}
}

func (s *exchangeService) Create(
	userId string,
	userEmail string,
	desiredBookID string,
	userBookIDs []string,
	latitude float64,
	longitude float64,
) (*models.ExchangeRequest, error) {
	offeredBooks := make([]models.OfferedBook, len(userBookIDs))
	for i, id := range userBookIDs {
		offeredBooks[i] = models.OfferedBook{BookID: id}
	}

	exchange := &models.ExchangeRequest{
		UserGoogleId:  userId,
		//UserEmail:     userEmail,
		DesiredBookID: desiredBookID,
		OfferedBooks:  offeredBooks,
		Status:        models.ExchangeRequestStatusActive,
		Latitude:      latitude,
		Longitude:     longitude,
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

func (s *exchangeService) GetAllWithStatus(userId string, status models.ExchangeRequestStatus) ([]*models.ExchangeRequest, error) {
	return s.repo.GetAllWithStatus(userId, status)
}

func (s *exchangeService) Get(id, userId string) (*models.ExchangeRequest, error) {
	return s.repo.Get(id, userId)
}

func (s *exchangeService) Delete(id string) error {
	r, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if r.Status == models.ExchangeRequestStatusCompleted {
		return models.ErrExchangeRequestCompleted
	}

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

	if r.Status == models.ExchangeRequestStatusCompleted {
		return matchingRequests, nil
	}

	if err := s.repo.Update(r); err != nil {
		return nil, err
	}

	for _, matchingReq := range matchingRequests {
		_, err := s.CreateMatch(r, matchingReq)
		if err != nil {
			return nil, err
		}
	}

	return matchingRequests, nil
}

func getOfferedBookIDs(offeredBooks []models.OfferedBook) []string {
	ids := make([]string, len(offeredBooks))
	for i, book := range offeredBooks {
		ids[i] = book.BookID
	}
	return ids
}

func (s *exchangeService) GetLocalizationAutocomplete(query string) ([]geo.Result, error) {
	return geo.GetLocalizationAutocomplete(query)
}
