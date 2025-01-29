package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID            string  `gorm:"primaryKey"`
	ISBN          uint    `json:"isbn"`
	Title         string  `json:"title"`
	Authors       string  `json:"authors"`
	Description   string  `json:"description"`
	ImageLink     string  `json:"thumbnail"`
	Genres        []Genre `gorm:"many2many:book_genres;constraint:OnDelete:CASCADE;"`
	Link          string
	PublishedDate string
	Pages         int
}
