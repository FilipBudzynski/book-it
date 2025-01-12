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

func (r *ExchangeRequestRepository) Get(id, userId string) (*models.ExchangeRequest, error) {
	exchange := &models.ExchangeRequest{}
	return exchange, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Where("user_google_id = ?", userId).
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

func (r *ExchangeRequestRepository) FindMatchingRequests(requestId, userId, desiredBookId string, offeredBooks []string) ([]*models.ExchangeRequest, error) {
	var matches []*models.ExchangeRequest
	err := r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Not("id = ?", requestId).                                                                        // Exclude user's own request
		Not("user_google_id = ?", userId).                                                               // Exclude user's own requests
		Where("desired_book_id IN ?", offeredBooks).                                                     // Their desired book is in your offered books
		Where("id IN (SELECT exchange_request_id FROM offered_books WHERE book_id = ?)", desiredBookId). // Your desired book is in their offered books
		Where("status = ?", models.ExchangeRequestStatusPending).                                        // Only match pending requests
		Find(&matches).Error
	if err != nil {
		return nil, err
	}
	return matches, nil
}
