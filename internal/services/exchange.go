package services

import (
	"github.com/FilipBudzynski/book_it/internal/models"
)

type ExchangeRequestRepository interface {
	Create(exchange *models.ExchangeRequest) error
}

type exchangeService struct {
	repo ExchangeRequestRepository
}

func NewExchangeService(r ExchangeRequestRepository) *exchangeService {
	return &exchangeService{
		repo: r,
	}
}

func (s *exchangeService) Create(userId, desiredBookID string, userBookIDs []string) error {
	offeredBooks := make([]models.OfferedBook, len(userBookIDs))
	for i, id := range userBookIDs {
		offeredBooks[i] = models.OfferedBook{BookId: id}
	}

	exchange := &models.ExchangeRequest{
		UserGoogleId:  userId,
		DesiredBookID: desiredBookID,
		OfferedBooks:  offeredBooks,
	}

	return s.repo.Create(exchange)
}
