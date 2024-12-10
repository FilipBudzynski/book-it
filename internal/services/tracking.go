package services

import (
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type TrackingService interface {
	// standard methods
	Create(progress *models.ReadingProgress) error
	Read(id string) (*models.ReadingProgress, error)
	Update(progress *models.ReadingProgress) error
	Delete(id string) error

	// custom methods
	GetByBookId(bookId string) (*models.ReadingProgress, error)
	GetDailyLog(bookId string, date time.Time) (*models.DailyReadingLog, error)
}

type trackingService struct {
	db *gorm.DB
}

func NewTrackingService(db *gorm.DB) TrackingService {
	return &trackingService{
		db: db,
	}
}

func (s *trackingService) Create(progress *models.ReadingProgress) error {
	return s.db.Create(progress).Error
}

func (s *trackingService) Read(id string) (*models.ReadingProgress, error) {
	readingProgress := &models.ReadingProgress{}
	err := s.db.First(readingProgress, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return readingProgress, nil
}

func (s *trackingService) Update(progress *models.ReadingProgress) error {
	return s.db.Save(progress).Error
}

func (s *trackingService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.ReadingProgress{}).Error
}

func (s *trackingService) GetByBookId(bookId string) (*models.ReadingProgress, error) {
	readingProgress := &models.ReadingProgress{}
	err := s.db.First(readingProgress, "book_id = ?", bookId).Error
	if err != nil {
		return nil, err
	}
	return readingProgress, nil
}

func (s *trackingService) GetDailyLog(bookId string, date time.Time) (*models.DailyReadingLog, error) {
	readingLog := &models.DailyReadingLog{}
	err := s.db.First(readingLog, "book_id = ? AND date = ?", bookId, date).Error
	if err != nil {
		return nil, err
	}
	return readingLog, nil
}
