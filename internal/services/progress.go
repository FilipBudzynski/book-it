package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
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
	if days > models.MaxDailyLogs {
		return models.ReadingProgress{}, models.ErrProgressMaxLogsExceeded
	}

	targetPages := int((totalPages + days - 1) / days)

	progressLogs := []models.DailyProgressLog{}
	for i := range days {
		logDate := startDateParsed.AddDate(0, 0, i)
		logTargetPages := CalculateTargetPages(totalPages, days-i)
		progressLog := &models.DailyProgressLog{
			Date:        logDate,
			TargetPages: logTargetPages,
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

func (s *progressService) RefreshTargetPagesForNewDay(progressId uint, date time.Time) error {
	return s.updateTargetPages(progressId, date, false)
}

func (s *progressService) UpdateTargetPages(progressId uint, updatedLogDate time.Time) error {
	return s.updateTargetPages(progressId, updatedLogDate, true)
}

func (s *progressService) updateTargetPages(progressId uint, referenceDate time.Time, includeEqualDate bool) error {
	progress, err := s.repo.GetById(strconv.FormatUint(uint64(progressId), 10))
	if err != nil {
		return err
	}

	original := *progress

	pagesLeft := progress.PagesLeft()
	progress.DailyTargetPages = CalculateTargetPages(
		pagesLeft,
		progress.DaysLeft(utils.TodaysDate())+1)

	if err := progress.Validate(); err != nil {
		return err
	}

	for i := range progress.DailyProgress {
		log := &progress.DailyProgress[i]
		if log.Date.Before(referenceDate) || (includeEqualDate && log.Date.Equal(referenceDate)) {
			continue
		}

		if log.Date.After(utils.TodaysDate()) {
			break
		}

		log.TargetPages = CalculateTargetPages(pagesLeft, log.DaysLeft(progress.EndDate))
		if err := log.Validate(); err != nil {
			return err
		}
	}

	if !progress.Equal(original) {
		return s.repo.Update(progress)
	}
	return nil
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
