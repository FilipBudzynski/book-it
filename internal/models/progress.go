package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

const MaxDailyLogs = 365

type ReadingProgress struct {
	gorm.Model
	UserBookID       uint   `gorm:"not null" form:"user-book-id"`
	BookTitle        string `form:"book-title"`
	StartDate        time.Time
	EndDate          time.Time
	TotalPages       int `form:"total-pages"`
	CurrentPage      int `form:"current-page"`
	DailyTargetPages int
	DailyProgress    []DailyProgressLog `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Completed        bool
}

var (
	ErrProgressCurrentPageGreaterThanTotal = errors.New("current page cannot be greater than total pages")
	ErrProgressCurrentPageNegative         = errors.New("current page cannot be negative")
	ErrProgressDailyTargetPagesNegative    = errors.New("daily target pages cannot be negative")
	ErrProgressDaysLeftNegative            = errors.New("days left cannot be negative")
	ErrProgressInvalidEndDate              = errors.New("end date must be after start date")
	ErrProgressLastDayNotFinished          = errors.New("this is the last day - finish reading today or change the end date to add more days")
	ErrProgressPagesLeftNegative           = errors.New("pages left cannot be negative")
	ErrProgressMaxLogsExceeded             = errors.New("max 365 logs per book")
)

func (r *ReadingProgress) Validate() error {
	if r.CurrentPage > r.TotalPages {
		return ErrProgressCurrentPageGreaterThanTotal
	}

	if r.CurrentPage < 0 {
		return ErrProgressCurrentPageNegative
	}
	if r.CurrentPage > r.TotalPages {
		return ErrProgressCurrentPageGreaterThanTotal
	}

	if r.DailyTargetPages < 0 {
		return ErrProgressDailyTargetPagesNegative
	}

	if r.StartDate.After(r.EndDate) {
		return ErrProgressInvalidEndDate
	}

	if r.PagesLeft() < 0 {
		return ErrProgressPagesLeftNegative
	}

	if r.IsCompleted() {
		r.Completed = true
	} else {
		r.Completed = false
	}

	return nil
}

func (r *ReadingProgress) DaysLeft(date time.Time) int {
	return int(r.EndDate.Sub(date).Hours() / 24)
}

func (r *ReadingProgress) IsFinishedByEndDate(date time.Time) error {
	daysLeft := r.DaysLeft(date)
	if daysLeft == 0 && r.PagesLeft() > 0 {
		return ErrProgressLastDayNotFinished
	}
	return nil
}

func (r *ReadingProgress) PagesLeft() int {
	return r.TotalPages - r.CurrentPage
}

func (r *ReadingProgress) IsCompleted() bool {
	return r.CurrentPage == r.TotalPages
}

func (r *ReadingProgress) IsFinishedOnLastLog(logDate time.Time) bool {
	return r.DaysLeft(logDate) != 0 || r.PagesLeft() < 0
}

func (r *ReadingProgress) GetLatestPositiveLog() *DailyProgressLog {
	for i := len(r.DailyProgress) - 1; i >= 0; i-- {
		if r.DailyProgress[i].PagesRead > 0 {
			return &r.DailyProgress[i]
		}
	}
	return nil
}

func (r *ReadingProgress) Equal(other ReadingProgress) bool {
	if other.DailyTargetPages != r.DailyTargetPages {
		return false
	}

	if len(other.DailyProgress) != len(r.DailyProgress) {
		return false
	}

	for i := range other.DailyProgress {
		otherLog := other.DailyProgress[i]
		log := r.DailyProgress[i]

		if otherLog.TargetPages != log.TargetPages {
			return false
		}
	}

	return true
}

func (r *ReadingProgress) GetLogForDate(date time.Time) *DailyProgressLog {
	for _, log := range r.DailyProgress {
		if log.Date.Equal(date) {
			return &log
		}
	}
	return nil
}
