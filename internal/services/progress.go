package services

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type progressService struct {
	db *gorm.DB
}

func NewProgressService(db *gorm.DB) handlers.ProgressService {
	return &progressService{db: db}
}

func (s *progressService) Create(progress *models.ReadingProgress) error {
	return s.db.Create(progress).Error
}

func (s *progressService) Get(id string) (*models.ReadingProgress, error) {
	readingProgress := &models.ReadingProgress{}
	err := s.db.First(readingProgress, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return readingProgress, nil
}

func (s *progressService) Update(progress *models.ReadingProgress) error {
	return s.db.Save(progress).Error
}

func (s *progressService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.ReadingProgress{}).Error
}

func (s *progressService) GetByUserBookId(id string) (*models.ReadingProgress, error) {
	readingProgress := &models.ReadingProgress{}
	err := s.db.Preload("DailyProgress").First(readingProgress, "user_book_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return readingProgress, nil
}
