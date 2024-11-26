package schemas

import (
	"gorm.io/gorm"
)

// Book is a schema representing a book received from external providers e.g. Google Books
type Book struct {
	gorm.Model
	ID            string
	ISBN          uint     `json:"isbn"`
	Title         string   `json:"title"`
	Authors       []string `json:"authors"`
	Description   string   `json:"description"`
	ImageLink     string   `json:"thumbnail"`
	Genre         string
	Link          string
	PublishedDate string
}
