package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
)

type ProgressRepository interface {
	Create(progress models.ReadingProgress) error
	GetById(id string) (*models.ReadingProgress, error)
	GetByUserBookId(userBookId string) (*models.ReadingProgress, error)
	Update(progress *models.ReadingProgress) error
	Delete(id string) error
	// logs methods
	// TODO: move to progressLogRepository
	GetLogById(id string) (*models.DailyProgressLog, error)
	UpdateLog(log *models.DailyProgressLog) error
}

type progressService struct {
	repo ProgressRepository
}

func NewProgressService(repo ProgressRepository) *progressService {
	return &progressService{repo: repo}
}

func (s *progressService) Create(bookId uint, totalPages int, bookTitle, startDate, endDate string) (models.ReadingProgress, error) {
	startDateParsed, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return models.ReadingProgress{}, err
	}

	endDateParsed, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return models.ReadingProgress{}, err
	}

	days := int(endDateParsed.Sub(startDateParsed).Hours()/24) + 1
	if days <= 0 {
		return models.ReadingProgress{}, models.ErrProgressInvalidEndDate
	}

	targetPages := int((totalPages + days - 1) / days)

	progressLogs := []models.DailyProgressLog{}
	for i := range days {
		progressLog := &models.DailyProgressLog{
			Date:        startDateParsed.AddDate(0, 0, i),
			TargetPages: targetPages,
			UserBookID:  bookId,
			TotalPages:  totalPages,
			Completed:   false,
		}
		if err := progressLog.Validate(); err != nil {
			return models.ReadingProgress{}, err
		}
		progressLogs = append(progressLogs, *progressLog)
	}

	progress := models.ReadingProgress{
		UserBookID:       bookId,
		BookTitle:        bookTitle,
		StartDate:        startDateParsed,
		EndDate:          endDateParsed,
		TotalPages:       totalPages,
		DailyTargetPages: targetPages,
		DailyProgress:    progressLogs,
		CurrentPage:      0,
		Completed:        false,
	}

	errs := errors.Join(
		progress.Validate(),
		s.repo.Create(progress),
	)

	return progress, errs
}

func (s *progressService) Get(id string) (*models.ReadingProgress, error) {
	return s.repo.GetById(id)
}

func (s *progressService) GetProgressAssosiatedWithLogId(logId string) (*models.ReadingProgress, error) {
	log, err := s.GetLog(logId)
	if err != nil {
		return nil, err
	}
	return s.Get(strconv.Itoa(int(log.ReadingProgressID)))
}

func (s *progressService) UpdateLog(id string, pagesRead int, comment string) (*models.DailyProgressLog, error) {
	log, err := s.repo.GetLogById(id)
	if err != nil {
		return nil, err
	}

	log.PagesRead = pagesRead
	log.Comment = comment

	if err := log.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateLog(log); err != nil {
		return nil, err
	}

	return log, nil
}

func (s *progressService) UpdateTargetPages(progressId uint) error {
	progress, err := s.repo.GetById(strconv.FormatUint(uint64(progressId), 10))
	if err != nil {
		return err
	}

	latestLog := progress.GetLatestPositiveLog()

	progress.DailyTargetPages = CalculateTargetPages(
		progress.PagesLeft(),
		progress.DaysLeft(latestLog.Date))

	if err := progress.Validate(); err != nil {
		return err
	}

	progress.UpdateLogTargetPagesFromDate(latestLog.Date)

	return s.repo.Update(progress)
}

func CalculateTargetPages(pagesLeft, daysLeft int) int {
	if pagesLeft < 0 {
		return -1
	}

	switch {
	case daysLeft == 0:
		return pagesLeft
	case daysLeft > 0:
		return int((pagesLeft + daysLeft - 1) / daysLeft)
	default:
		return -1
	}
}

func (s *progressService) GetLog(id string) (*models.DailyProgressLog, error) {
	return s.repo.GetLogById(id)
}

func (s *progressService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *progressService) GetByUserBookId(id string) (*models.ReadingProgress, error) {
	return s.repo.GetByUserBookId(id)
}
