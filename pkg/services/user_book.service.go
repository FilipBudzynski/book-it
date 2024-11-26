package services

import (
	"github.com/FilipBudzynski/book_it/pkg/models"
	"gorm.io/gorm"
)

type userBookService struct {
	db *gorm.DB
}

func NewUserBookService(db *gorm.DB) *userBookService {
	return &userBookService{
		db: db,
	}
}

func (s *userBookService) Create(userId, bookId string) error {
	userBook := &models.UserBook{
		UserGoogleId: userId,
		BookId:       bookId,
		Status:       models.BookStatusNotStarted,
	}
	return s.db.Create(userBook).Error
}

func (s *userBookService) Update(userBook *models.UserBook) error {
	return s.db.Save(userBook).Error
}

func (s *userBookService) GetAll() ([]models.UserBook, error) {
	var userBooks []models.UserBook
	if err := s.db.Find(&userBooks).Error; err != nil {
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
