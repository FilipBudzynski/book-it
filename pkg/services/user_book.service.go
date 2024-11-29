package services

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/pkg/handlers"
	"github.com/FilipBudzynski/book_it/pkg/models"
	"github.com/FilipBudzynski/book_it/pkg/schemas"
	"gorm.io/gorm"
)

type userBookService struct {
	db          *gorm.DB
	bookService handlers.BookService
}

func NewUserBookService(db *gorm.DB, bookSerivce handlers.BookService) *userBookService {
	return &userBookService{
		db:          db,
		bookService: bookSerivce,
	}
}

func (s *userBookService) Create(userID, bookID string) error {
	userBook := &models.UserBook{
		UserGoogleId: userID,
		BookID:       bookID,
		Status:       models.BookStatusNotStarted,
	}
	return s.db.Create(userBook).Error
}

func (s *userBookService) Update(userBook *models.UserBook) error {
	return s.db.Save(userBook).Error
}

func (s *userBookService) Delete(id string) error {
	return s.db.Where("book_id = ?", id).Delete(&models.UserBook{}).Error
}

func (s *userBookService) GetUserBooks(userId string) ([]schemas.Book, error) {
	userBooks, err := s.GetAll(userId)
	if err != nil {
		return nil, err
	}

	var books []schemas.Book
	for _, userBook := range userBooks {
		book, err := s.bookService.GetByID(userBook.BookID)
		if err != nil {
			return nil, fmt.Errorf("Failed to get book with id: %s from the provider, err: %v", userBook.BookID, err)
		}
		books = append(books, book)
	}

	return books, nil
}

func (s *userBookService) GetAll(userId string) ([]models.UserBook, error) {
	var userBooks []models.UserBook
	if err := s.db.Where("user_google_id = ?", userId).Find(&userBooks).Error; err != nil {
		return nil, err
	}

	return userBooks, nil
}

func (s *userBookService) GetById(id string) (*models.UserBook, error) {
	var userBook models.UserBook
	if err := s.db.First(&userBook, id).Error; err != nil {
		return nil, err
	}
	return &userBook, nil
}
