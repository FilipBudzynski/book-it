package services

import (
	"github.com/FilipBudzynski/book_it/internal/models"
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

func (s *userBookService) GetAll(userId string) ([]*models.UserBook, error) {
	var userBooks []*models.UserBook
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
