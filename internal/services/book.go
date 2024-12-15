package services

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

// BookProvider is used to communicate with the external API or Database
// in order to retreive response and parse it into models.Book struct
type BookProvider interface {
	GetBook(id string) (*models.Book, error)
	GetBooksByQuery(query string, limit int) ([]*models.Book, error)
	// used to change the limit of query results
	WithLimit(limit int) BookProvider
}

type bookService struct {
	provider BookProvider
	db       *gorm.DB
}

func NewBookService(bookProvider BookProvider, db *gorm.DB) handlers.BookService {
	return &bookService{
		provider: bookProvider,
		db:       db,
	}
}

// get fetches the first book by isbn
func (s *bookService) get(id string) (*models.Book, error) {
	var book models.Book
	err := s.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (s *bookService) Delete(userID, bookID string) error {
	return s.db.Where("id = ?", bookID).Where("user_google_id = ?", userID).Delete(&models.Book{}).Error
}

func (s *bookService) Create(book *models.Book) error {
	return s.db.Create(book).Error
}

func (s *bookService) Get(id string) (*models.Book, error) {
	book := &models.Book{}
	err := s.db.First(book, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return book, nil
}

// GetByID fetches the book from Database
// if no book is found, it fetches the book from the provider and saves it to the database
func (s *bookService) GetByID(bookId string) (*models.Book, error) {
	if book, _ := s.get(bookId); book != nil {
		return book, nil
	}

	book, err := s.provider.GetBook(bookId)
	if err != nil {
		return nil, err
	}

	if err := s.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

// GetByQuery returns maxResults number of books by title from external api (provider)
// if no book is found in database, it saves it to the database
func (s *bookService) GetByQuery(query string, limit int) ([]*models.Book, error) {
	books, err := s.provider.GetBooksByQuery(query, limit)
	if err != nil {
		return nil, err
	}

	// store book in database if not found
	for _, book := range books {
		if dbBook, _ := s.Get(book.ID); dbBook != nil {
			continue
		}
		if err := s.Create(book); err != nil {
			return nil, err
		}
	}
	return books, nil
}
