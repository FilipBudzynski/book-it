package models

import "gorm.io/gorm"

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
