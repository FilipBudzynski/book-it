package services

import (
	"math/rand/v2"
	"slices"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
)

const MaxRecommendationsResults = 30

type BookRepository interface {
	Create(book *models.Book) error
	Get(id string) (*models.Book, error)
	Delete(id string) error
	GetByGenre(genre string) ([]*models.Book, error)
}

type bookService struct {
	provider handlers.BookProvider
	repo     BookRepository
}

func NewBookService(repo BookRepository) handlers.BookService {
	return &bookService{
		repo: repo,
	}
}

func (s *bookService) WithProvider(provider handlers.BookProvider) handlers.BookService {
	s.provider = provider
	return s
}

func (s *bookService) Provider() handlers.BookProvider {
	return s.provider
}

func (s *bookService) Delete(bookID string) error {
	return s.repo.Delete(bookID)
}

func (s *bookService) Create(book *models.Book) error {
	return s.repo.Create(book)
}

func (s *bookService) GetByID(bookId string) (*models.Book, error) {
	book, err := s.repo.Get(bookId)
	if err == nil && book != nil {
		return book, nil
	}

	book, err = s.provider.GetBook(bookId)
	if err != nil {
		return nil, err
	}

	if err := s.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) GetByQuery(query string, queryType handlers.QueryType, page int) ([]*models.Book, error) {
	books, err := s.provider.GetBooksByQuery(query, queryType, 40, page)
	if err != nil {
		return nil, err
	}

	// store book in database if not found
	for _, book := range books {
		if dbBook, _ := s.repo.Get(book.ID); dbBook != nil {
			continue
		}
		if err := s.Create(book); err != nil {
			return nil, err
		}
	}
	return books, nil
}

func (s *bookService) FetchReccomendations(genres []models.Genre, userBooks []*models.UserBook) ([]*models.Book, error) {
	userBookIDs := []string{}
	for _, userBook := range userBooks {
		userBookIDs = append(userBookIDs, userBook.Book.ID)
	}

	providerBooks := []*models.Book{}
	for _, genre := range genres {
		genreBooks, err := s.provider.GetBooksByGenre(genre.Name)
		if err != nil {
			return nil, err
		}
		providerBooks = append(providerBooks, genreBooks...)
	}

	for _, book := range providerBooks {
		if dbBook, _ := s.repo.Get(book.ID); dbBook != nil {
			continue
		}
		if err := s.Create(book); err != nil {
			return nil, err
		}
	}

	resultBooks := []*models.Book{}
	for _, book := range providerBooks {
		if slices.Contains(userBookIDs, book.ID) {
			continue
		}
		resultBooks = append(resultBooks, book)
	}

	rand.Shuffle(len(resultBooks), func(i, j int) {
		resultBooks[i], resultBooks[j] = resultBooks[j], resultBooks[i]
	})

	if len(resultBooks) < MaxRecommendationsResults {
		return resultBooks, nil
	}
	return resultBooks[:MaxRecommendationsResults], nil
}
