package services

import (
	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
)

func NewBookService(bookProvider BookProvider) handlers.BookService {
	return &bookService{
		provider: bookProvider,
	}
}

// BookProvider is used to communicate with the external API or Database
// in order to retreive response and parse it into models.Book struct
type BookProvider interface {
	GetBook(id string) (schemas.Book, error)
	GetBooksByQuery(query string, limit int) ([]schemas.Book, error)
	// used to change the limit of query results
	WithLimit(limit int) BookProvider
}

type bookService struct {
	provider BookProvider
}

func (s *bookService) GetByID(bookId string) (schemas.Book, error) {
	book, err := s.provider.GetBook(bookId)
	if err != nil {
		return schemas.Book{}, err
	}

	return book, nil
}

func (s *bookService) GetByQuery(query string, limit int) ([]schemas.Book, error) {
	books, err := s.provider.GetBooksByQuery(query, limit)
	if err != nil {
		return nil, err
	}

	return books, nil
}
