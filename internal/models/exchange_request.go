package models

import "gorm.io/gorm"

type ExchangeRequest struct {
	gorm.Model
	UserGoogleId  string `form:"user_id"`
	User          User   `gorm:"foreignKey:UserGoogleId"`
	DesiredBookID string `gorm:"foreignKey:BookID" form:"book_id"`
	DesiredBook   Book   `gorm:"foreignKey:DesiredBookID;constraint:OnDelete:SET NULL"`
	OfferedBooks  []OfferedBook
}

type OfferedBook struct {
	gorm.Model
	ExchangeRequestID uint
	BookId            string `gorm:"foreignKey:BookID" form:"book_id"`
}
