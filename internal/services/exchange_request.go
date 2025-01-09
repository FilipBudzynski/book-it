package services

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/internal/models"
)

type ExchangeRequestRepository interface {
	Create(exchange *models.ExchangeRequest) error
	Get(id string) (*models.ExchangeRequest, error)
	GetAll(userId string) ([]*models.ExchangeRequest, error)
	Delete(id string) error
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
	return s.repo.Get(fmt.Sprintf("%d", exchange.ID))
}

func (s *exchangeService) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	return s.repo.GetAll(userId)
}

func (s *exchangeService) Get(id string) (*models.ExchangeRequest, error) {
	return s.repo.Get(id)
}

func (s *exchangeService) Delete(id string) error {
	return s.repo.Delete(id)
}
