package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type DailyProgressLog struct {
	gorm.Model
	ID                uint `gorm:"primaryKey"`
	ReadingProgressID uint `gorm:"not null"` // Reading progress foreign key
	UserBookID        uint // Denormalized
	Date              time.Time
	PagesRead         int `form:"pages-read"`  // Pages read on this date
	TotalPages        int `form:"total-pages"` // Denormalized
	TargetPages       int
	Completed         bool // Whether the day's target was met
	Comment           string
}

var (
	ErrProgressLogTrackingEndsBeforeStart   = errors.New("end date must be after start date")
	ErrProgressLogPagesReadNotSpecified     = errors.New("pages read must be a positive number")
	ErrProgressLogPagesReadGreaterThanTotal = errors.New("pages read cannot be greater than total pages")
	ErrProgressLogUpdateDateInFuture        = errors.New("cannot update a log in the future")
)

func (d *DailyProgressLog) Validate() error {
	if d.PagesRead < 0 {
		return ErrProgressLogPagesReadNotSpecified
	}

	if d.PagesRead >= d.TargetPages {
		d.Completed = true
	} else {
		d.Completed = false
	}

	return nil
}

func (d *DailyProgressLog) AfterSave(db *gorm.DB) error {
	var readingProgress ReadingProgress
	err := db.First(&readingProgress, d.ReadingProgressID).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	var totalPagesRead int
	if err := db.Model(&DailyProgressLog{}).
		Where("reading_progress_id = ?", d.ReadingProgressID).
		Select("SUM(pages_read)").
		Row().
		Scan(&totalPagesRead); err != nil {
		return err
	}

	readingProgress.CurrentPage = int(totalPagesRead)
	if err := readingProgress.Validate(); err != nil {
		return err
	}

	if err := db.Save(&readingProgress).Error; err != nil {
		return err
	}

	return nil
}

func (d *DailyProgressLog) DaysLeft(endDate time.Time) int {
	return int(endDate.Sub(d.Date).Hours()/24) + 1 // +1 because endDate is exclusive
}

func (d *DailyProgressLog) IsEmptyOrOverdue(date time.Time) bool {
	if d.PagesRead == 0 || d.Date.Before(date) {
		return true
	}
	return false
}
