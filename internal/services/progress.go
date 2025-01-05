package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
)

var ErrDailyTargetPagesNegative = errors.New("something went wrong with the count of daily logs")

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

func NewProgressService(repo ProgressRepository) handlers.ProgressService {
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

	days := int(endDateParsed.Sub(startDateParsed).Hours() / 24)
	if days <= 0 {
		return models.ReadingProgress{}, errors.New("end date must be after start date")
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

func (s *progressService) GetProgressByAssosiatedLogId(id string) (*models.ReadingProgress, error) {
	log, err := s.GetLog(id)
	if err != nil {
		return nil, err
	}
	return s.Get(strconv.Itoa(int(log.ReadingProgressID)))
}

func (s *progressService) UpdateLogPagesRead(id, pagesReadString string) error {
	if pagesReadString == "" {
		return models.ErrProgressLogPagesReadNotSpecified
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

	if errs := errors.Join(
		log.Validate(),
		s.repo.UpdateLog(log),
		s.updateTargetPages(log.ReadingProgressID, log.Date),
	); errs != nil {
		return errs
	}

	return nil
}

func (s *progressService) updateTargetPages(progressId uint, logDate time.Time) error {
	progress, err := s.repo.GetById(strconv.FormatUint(uint64(progressId), 10))
	if err != nil {
		return err
	}

	pagesLeft := progress.TotalPages - progress.CurrentPage
	daysLeft := int(progress.EndDate.Sub(logDate).Hours()/24) - 1
	var targetPages int

	switch {
	case daysLeft == 0:
		targetPages = pagesLeft
	case daysLeft > 0:
		targetPages = int((pagesLeft + daysLeft - 1) / daysLeft)
	case daysLeft < 0:
		return models.ErrProgressDaysLeftNegative
	}

	progress.DailyTargetPages = targetPages
	for i := range progress.DailyProgress {
		progress.DailyProgress[i].TargetPages = targetPages
	}

	return errors.Join(
		progress.Validate(),
		s.repo.Update(progress),
		s.checkComplete(daysLeft, pagesLeft),
	)
}

func (s *progressService) checkComplete(daysLeft, pagesLeft int) error {
	if daysLeft == 0 && pagesLeft > 0 {
		return models.ErrProgressPastEndDate
	}
	return nil
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
