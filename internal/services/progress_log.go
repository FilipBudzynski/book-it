package services

import (
	"time"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type progressLogService struct {
	db *gorm.DB
}

func NewProgressLogService(db *gorm.DB) handlers.ProgressLogService {
	return &progressLogService{db: db}
}

func (s *progressLogService) Create(progressId uint, target int, date time.Time) (*models.DailyProgressLog, error) {
	readingLog := &models.DailyProgressLog{
		ReadingProgressID: progressId,
		Date:              date,
		TargetPages:       target,
		Completed:         false,
	}

	err := s.db.Create(readingLog).Error
	if err != nil {
		return nil, err
	}
	return readingLog, nil
}

func (s *progressLogService) Update(progressLog *models.DailyProgressLog) error {
	return s.db.Save(progressLog).Error
}

func (s *progressLogService) Get(id string) (*models.DailyProgressLog, error) {
	progressLog := &models.DailyProgressLog{}
	err := s.db.First(&progressLog, id).Error
	if err != nil {
		return nil, err
	}
	return progressLog, nil
}

func (s *progressLogService) GetAll(progressId string) ([]models.DailyProgressLog, error) {
	var progressLogs []models.DailyProgressLog
	err := s.db.Find(&progressLogs, "reading_progress_id = ?", progressId).Error
	if err != nil {
		return nil, err
	}
	return progressLogs, nil
}

func (s *progressLogService) Delete(id string) error {
	return s.db.Where("id = ?", id).Delete(&models.DailyProgressLog{}).Error
}
