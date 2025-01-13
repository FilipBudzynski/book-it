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
	FindMatchingRequests(userId, requestId, desiredBookId string, offeredBooks []string, existingMatches []models.ExchangeMatch) ([]models.ExchangeMatch, error)
}

type exchangeService struct {
	repo ExchangeRequestRepository
}

func NewExchangeService(r ExchangeRequestRepository) *exchangeService {
	return &exchangeService{
		repo: r,
	}
}

func (s *exchangeService) Create(userId, desiredBookID string, userBookIDs []string) (*models.ExchangeRequest, error) {
	offeredBooks := make([]models.OfferedBook, len(userBookIDs))
	for i, id := range userBookIDs {
		offeredBooks[i] = models.OfferedBook{BookId: id}
	}

	exchange := &models.ExchangeRequest{
		UserGoogleId:  userId,
		DesiredBookID: desiredBookID,
		OfferedBooks:  offeredBooks,
		Status:        models.ExchangeRequestStatusPending,
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

func (s *exchangeService) FindMatchingRequests(requestId, userId string) (*models.ExchangeRequest, error) {
	r, err := s.Get(requestId, userId)
	if err != nil {
		return nil, err
	}

	matches, err := s.repo.FindMatchingRequests(
		requestId,
		userId,
		r.DesiredBookID,
		getOfferedBookIDs(r.OfferedBooks),
		r.Matches,
	)
	if err != nil {
		return nil, err
	}

	r.Matches = append(r.Matches, matches...)
    if len(r.Matches) == 0 {
        r.Status = models.ExchangeRequestStatusPending
    }

	if err := s.repo.Update(r); err != nil {
		return nil, err
	}

	return r, nil
}

// Helper function to extract BookIDs from OfferedBooks
func getOfferedBookIDs(offeredBooks []models.OfferedBook) []string {
	ids := make([]string, len(offeredBooks))
	for i, book := range offeredBooks {
		ids[i] = book.BookId
	}
	return ids
}
