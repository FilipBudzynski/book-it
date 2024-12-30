package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var MigrateModels = []any{
	&User{},
	&UserBook{},
	&Book{},
	&ReadingProgress{},
	&DailyProgressLog{},
}

type bookTrackingStatus string

const (
	BookStatusNotStarted bookTrackingStatus = "NotStarted"
	BookStatusStarted    bookTrackingStatus = "Started"
	BookStatusFinished   bookTrackingStatus = "Finished"
)

type User struct {
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	GoogleId string `gorm:"primaryKey" json:"google_id"`
	Books    []UserBook
	gorm.Model
}

type Book struct {
	gorm.Model
	ID            string `gorm:"primaryKey"`
	ISBN          uint   `json:"isbn"`
	Title         string `json:"title"`
	Authors       string `json:"authors"`
	Description   string `json:"description"`
	ImageLink     string `json:"thumbnail"`
	Genre         string
	Link          string
	PublishedDate string
	InBookShelff  bool
	Pages         int
}

// UserBook model is an abstraction for link between user and a book
// It bounds book to a user providing more information about user interactions with a book
type UserBook struct {
	gorm.Model
	IsTracked       bool   `gorm:"not null;default:false"`
	UserGoogleId    string `gorm:"not null"` // foreignKey
	BookID          string `gorm:"not null"`
	Book            Book
	ReadingProgress *ReadingProgress `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type ReadingProgress struct {
	gorm.Model
	UserBookID       uint `gorm:"not null" form:"user-book-id"` // Reference to the user book being read
	StartDate        time.Time
	EndDate          time.Time
	TotalPages       int                `form:"total-pages"`  // Total pages in the book
	CurrentPage      int                `form:"current-page"` // Curent page the user is on
	DailyTargetPages int                // Pages that need to be read a day to finish the book on time
	DailyProgress    []DailyProgressLog // Which days user have not read the book
	Completed        bool               // Whether the book is finished
}

func (r *ReadingProgress) Validate() error {
	if r.CurrentPage > r.TotalPages {
		return errors.New("Current page cannot be greater than total pages")
	}

	if r.CurrentPage < 0 {
		return errors.New("Current page cannot be negative")
	}

	if r.DailyTargetPages < 0 {
		return errors.New("Daily target pages cannot be negative")
	}

	if r.StartDate.After(r.EndDate) {
		return errors.New("End date must be after start date")
	}

	return nil
}

type DailyProgressLog struct {
	gorm.Model
	ReadingProgressID uint `gorm:"not null"` // Reading progress foreign key
	UserBookID        uint // Denormalized
	Date              time.Time
	PagesRead         int `form:"pages-read"`  // Pages read on this date
	TotalPages        int `form:"total-pages"` // Denormalized
	TargetPages       int
	Completed         bool // Whether the day's target was met
}

var (
	ErrTrackingEndsBeforeStart   = errors.New("End date must be after start date")
	ErrPagesReadNotSpecified     = errors.New("Pages read must be a positive number")
	ErrPagesReadGreaterThanTotal = errors.New("Pages read cannot be greater than total pages")
)

func (d *DailyProgressLog) Validate() error {
	if d.PagesRead < 0 {
		return ErrPagesReadNotSpecified
	}

	if d.PagesRead > d.TotalPages {
		return ErrPagesReadGreaterThanTotal
	}

	if d.PagesRead >= d.TargetPages {
		d.Completed = true
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

	if err := db.Save(&readingProgress).Error; err != nil {
		return err
	}
	return nil
}

type Bookshelf struct {
	gorm.Model
	Name         string
	UserGoogleId string `gorm:"not null"` // foreignKey
	Books        []UserBook
}
