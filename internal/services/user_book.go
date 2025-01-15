package services

import (
	"github.com/FilipBudzynski/book_it/internal/models"
)

type UserBookRepository interface {
	Create(userBook *models.UserBook) error
	GetAllUserBooks(userId string) ([]*models.UserBook, error)
	Get(id string) (*models.UserBook, error)
	Update(userBook *models.UserBook) error
	Delete(id string) error
	DeleteWhereBookId(bookId string) error
	Search(userId, search string) ([]*models.UserBook, error)
}

type userBookService struct {
	repo         UserBookRepository
	exchangeRepo ExchangeRequestRepository
}

func NewUserBookService(userBookRepo UserBookRepository, exchangeRepo ExchangeRequestRepository) *userBookService {
	return &userBookService{
		repo:         userBookRepo,
		exchangeRepo: exchangeRepo,
	}
}

func (s *userBookService) Create(userID, bookID string) error {
	userBook := &models.UserBook{
		UserGoogleId: userID,
		BookID:       bookID,
	}

	return s.repo.Create(userBook)
}

func (s *userBookService) Update(userBook *models.UserBook) error {
	return s.repo.Update(userBook)
}

func (s *userBookService) Get(id string) (*models.UserBook, error) {
	return s.repo.Get(id)
}

func (s *userBookService) GetAll(userId string) ([]*models.UserBook, error) {
	return s.repo.GetAllUserBooks(userId)
}

func (s *userBookService) Delete(id string) error {
	userBook, err := s.repo.Get(id)
	if err != nil {
		return err
	}

	activeRequests, err := s.exchangeRepo.GetActiveExchangeRequestsByBookID(userBook.BookID)
	if err != nil {
		return err
	}

	if len(activeRequests) > 0 {
		return models.ErrUserBookInActiveExchangeRequest
	}
	return s.repo.Delete(id)
}

func (s *userBookService) DeleteByBookId(bookId string) error {
	return s.repo.DeleteWhereBookId(bookId)
}

func (s *userBookService) Search(userId, query string) ([]*models.UserBook, error) {
	return s.repo.Search(userId, query)
}
