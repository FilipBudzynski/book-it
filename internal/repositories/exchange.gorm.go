package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type ExchangeRequestRepository struct {
	db *gorm.DB
}

func NewExchangeRequestRepository(db *gorm.DB) *ExchangeRequestRepository {
	return &ExchangeRequestRepository{
		db: db,
	}
}

func (r *ExchangeRequestRepository) Create(exchange *models.ExchangeRequest) error {
	return r.db.Create(exchange).Error
}
