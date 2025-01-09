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

func (r *ExchangeRequestRepository) Get(id string) (*models.ExchangeRequest, error) {
	exchange := &models.ExchangeRequest{}
	return exchange, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		First(&exchange, id).Error
}

func (r *ExchangeRequestRepository) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	exchanges := []*models.ExchangeRequest{}
	return exchanges, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Where("user_google_id = ?", userId).
		Find(&exchanges).Error
}

func (r *ExchangeRequestRepository) Delete(id string) error {
	return r.db.Delete(&models.ExchangeRequest{}, id).Error
}
