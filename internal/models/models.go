package models

import (
	"time"

	"gorm.io/gorm"
)

var MigrateModels = []any{
	&User{},
	&UserBook{},
	&Book{},
	&ReadingProgress{},
}

type userBookStatus string

const (
	BookStatusNotStarted userBookStatus = "NotStarted"
	BookStatusStarted    userBookStatus = "Started"
	BookStatusFinished   userBookStatus = "Finished"
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
}

// UserBook model is an abstraction for link between user and a book
// It bounds book to a user providing more information about user interactions with a book
type UserBook struct {
	gorm.Model
	Status       userBookStatus `gorm:"not null"`
	UserGoogleId string         `gorm:"not null"` // foreignKey
	BookID       string         `gorm:"not null"`
}

type ReadingProgress struct {
	gorm.Model
	UserBookID       string `gorm:"not null"` // Reference to the user book being read
	StartDate        time.Time
	EndDate          time.Time
	TotalPages       int               // Total pages in the book
	CurrentPage      int               // Curent page the user is on
	DailyTargetPages int               // Pages that need to be read a day to finish the book on time
	DailyProgress    []DailyReadingLog // Which days user have not read the book
	Completed        bool              // Whether the book is finished
}

type DailyReadingLog struct {
	gorm.Model
	ReadingProgressID uint `gorm:"not null"` // Reading progress foreign key
	Date              time.Time
	PagesRead         int  // Pages read on this date
	Completed         bool // Whether the day's target was met
}

type Bookshelf struct {
	gorm.Model
	Name         string
	UserGoogleId string `gorm:"not null"` // foreignKey
	Books        []UserBook
}
