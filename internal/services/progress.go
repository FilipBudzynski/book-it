package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
)

var (
	ErrorTrackingEndsBeforeStart = errors.New("End date must be after start date")
	ErrPagesReadNotSpecified     = errors.New("Pages read must be a positive number")
)

type ProgressRepository interface {
	Create(progress models.ReadingProgress) error
	GetById(id string, preloads ...string) (*models.ReadingProgress, error)
	GetByUserBookId(userBookId string, preloads ...string) (*models.ReadingProgress, error)
	Delete(id string) error
	// logs methods
	// TODO: move to progressLogRepository
	GetLogById(id string) (*models.DailyProgressLog, error)
	UpdateLog(log *models.DailyProgressLog) error
}

type progressService struct {
	repo ProgressRepository
}

func NewProgressService(repo ProgressRepository) handlers.ProgressService {
	return &progressService{repo: repo}
}

func (s *progressService) Create(bookId uint, totalPages int, startDate, endDate string) (models.ReadingProgress, error) {
	startDateParsed, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return models.ReadingProgress{}, err
	}

	endDateParsed, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return models.ReadingProgress{}, err
	}

	days := int(endDateParsed.Sub(startDateParsed).Hours() / 24)
	if days == 0 {
		return models.ReadingProgress{}, errors.New("End date must be after start date")
	}

	targetPages := int(totalPages / days)

	progressLogs := []models.DailyProgressLog{}
	for i := range days {
		progressLog := &models.DailyProgressLog{
			Date:        startDateParsed.AddDate(0, 0, i),
			TargetPages: targetPages,
			UserBookID:  bookId,
			TotalPages:  totalPages,
			Completed:   false,
		}
		progressLogs = append(progressLogs, *progressLog)
	}

	progress := models.ReadingProgress{
		UserBookID:       bookId,
		StartDate:        startDateParsed,
		EndDate:          endDateParsed,
		TotalPages:       totalPages,
		DailyTargetPages: targetPages,
		DailyProgress:    progressLogs,
		CurrentPage:      0,
		Completed:        false,
	}

	err = errors.Join(
		progress.Validate(),
		s.repo.Create(progress),
	)

	return progress, err
}

func (s *progressService) Get(id string) (*models.ReadingProgress, error) {
	return s.repo.GetById(id, "DailyProgress")
}

func (s *progressService) UpdateLogPagesRead(id, pagesReadString string) error {
	if pagesReadString == "" {
		return ErrPagesReadNotSpecified
	}

	pagesRead, err := strconv.Atoi(pagesReadString)
	if err != nil {
		return err
	}

	log, err := s.repo.GetLogById(id)
	if err != nil {
		return err
	}

	log.PagesRead = pagesRead

	return errors.Join(
		log.Validate(),
		s.repo.UpdateLog(log),
	)
}

func (s *progressService) GetLog(id string) (*models.DailyProgressLog, error) {
	return s.repo.GetLogById(id)
}

func (s *progressService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *progressService) GetByUserBookId(id string) (*models.ReadingProgress, error) {
	return s.repo.GetByUserBookId(id, "DailyProgress")
}
