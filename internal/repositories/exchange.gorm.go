package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		Preload("Matches.MatchRequest.DesiredBook").
		Where("user_google_id = ?", userId).
		First(&exchange, id).Error
}

func (r *ExchangeRequestRepository) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	exchanges := []*models.ExchangeRequest{}
	return exchanges, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Preload("Matches.MatchRequest.DesiredBook").
		Where("user_google_id = ?", userId).
		Find(&exchanges).Error
}

func (r *ExchangeRequestRepository) Delete(id string) error {
	return r.db.Select(clause.Associations).Delete(&models.ExchangeRequest{}, id).Error
}

func (r *ExchangeRequestRepository) DeleteMatchesForRequest(requestId string) error {
	if err := r.db.Where("exchange_request_id = ? OR match_request_id = ?", requestId, requestId).Delete(&models.ExchangeMatch{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ExchangeRequestRepository) Update(exchange *models.ExchangeRequest) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(exchange).Error
}

func (r *ExchangeRequestRepository) FindMatchingRequests(requestId, userId, desiredBookId string, offeredBooks []string) ([]*models.ExchangeRequest, error) {
	matches := []*models.ExchangeRequest{}
	err := r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Not("id = ?", requestId).                                                                        // Exclude user's own request
		Not("user_google_id = ?", userId).                                                               // Exclude user's own requests
		Where("desired_book_id IN ?", offeredBooks).                                                     // Their desired book is in your offered books
		Where("id IN (SELECT exchange_request_id FROM offered_books WHERE book_id = ?)", desiredBookId). // Your desired book is in their offered books
		Not("status = ?", models.ExchangeRequestStatusAccepted).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *ExchangeRequestRepository) CreateMatch(requestId, matchId uint, status models.ExchangeRequestStatus) (*models.ExchangeMatch, error) {
	exchangeMatch := &models.ExchangeMatch{
		ExchangeRequestID: requestId,
		MatchRequestID:    matchId,
		Status:            status,
	}
	return exchangeMatch, r.db.Create(exchangeMatch).Error
}

func (r *ExchangeRequestRepository) GetMatch(userReqId, otherReqId uint) (*models.ExchangeMatch, error) {
	exchangeMatch := &models.ExchangeMatch{}
	return exchangeMatch, r.db.First(&exchangeMatch, "exchange_request_id = ? AND match_request_id = ?", userReqId, otherReqId).Error
}

func (r *ExchangeRequestRepository) UpdateMatch(match *models.ExchangeMatch) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(match).Error
}
