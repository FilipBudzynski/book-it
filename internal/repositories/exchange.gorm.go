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

func (r *ExchangeRequestRepository) GetByID(id string) (*models.ExchangeRequest, error) {
	exchange := &models.ExchangeRequest{}
	return exchange, r.db.Preload("DesiredBook").Preload("OfferedBooks.Book").First(&exchange, id).Error
}

func (r *ExchangeRequestRepository) Get(id, userId string) (*models.ExchangeRequest, error) {
	exchange := &models.ExchangeRequest{}
	err := r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Where("user_google_id = ?", userId).
		First(&exchange, id).Error
	if err != nil {
		return nil, err
	}

	// Load Matches where this request is either the origin or the target
	var matches []models.ExchangeMatch
	err = r.db.Where("exchange_request_id = ? OR matched_exchange_request_id = ?", exchange.ID, exchange.ID).
		Preload("Request.DesiredBook").
		Preload("MatchedExchangeRequest.DesiredBook").
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	// Assign matches to the Matches field
	exchange.Matches = matches
	return exchange, nil
}

func (r *ExchangeRequestRepository) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	exchanges := []*models.ExchangeRequest{}
	return exchanges, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Preload("Matches.MatchedExchangeRequest.DesiredBook").
		Where("user_google_id = ?", userId).
		Find(&exchanges).Error
}

func (r *ExchangeRequestRepository) GetAllWithStatus(userId string, status models.ExchangeRequestStatus) ([]*models.ExchangeRequest, error) {
	exchanges := []*models.ExchangeRequest{}
	return exchanges, r.db.Preload("DesiredBook").
		Preload("OfferedBooks.Book").
		Preload("Matches.MatchedExchangeRequest.DesiredBook").
		Where("user_google_id = ?", userId).
		Where("status = ?", status.String()).
		Find(&exchanges).Error
}

func (r *ExchangeRequestRepository) Delete(id string) error {
	return r.db.Delete(&models.ExchangeRequest{}, id).Error
}

func (r *ExchangeRequestRepository) DeleteMatchesForRequest(requestId string) error {
	if err := r.db.Where("exchange_request_id = ? OR matched_exchange_request_id = ?", requestId, requestId).Delete(&models.ExchangeMatch{}).Error; err != nil {
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
		Not("status = ?", models.ExchangeRequestStatusCompleted).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *ExchangeRequestRepository) CreateMatch(match *models.ExchangeMatch) error {
	if match.ExchangeRequestID > match.MatchedExchangeRequestID {
		match.ExchangeRequestID, match.MatchedExchangeRequestID = match.MatchedExchangeRequestID, match.ExchangeRequestID
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "exchange_request_id"}, {Name: "matched_exchange_request_id"}},
		DoNothing: true,
	}).Create(match).Error
}

func (r *ExchangeRequestRepository) GetMatch(matchID, requestID string) (*models.ExchangeMatch, error) {
	matches, err := r.getMatches(requestID, matchID)
	return matches[0], err
}

func (r *ExchangeRequestRepository) GetAllMatches(requestId string) ([]*models.ExchangeMatch, error) {
	return r.getMatches(requestId, "")
}

func (r *ExchangeRequestRepository) getMatches(requestId string, matchID string) ([]*models.ExchangeMatch, error) {
	matches := []*models.ExchangeMatch{}
	query := r.db.Preload("Request").
		Preload("Request.DesiredBook").
		Preload("MatchedExchangeRequest").
		Preload("MatchedExchangeRequest.DesiredBook").
		Where("exchange_request_id = ? OR matched_exchange_request_id = ?", requestId, requestId)

	if matchID != "" {
		return matches, query.Where("id = ?", matchID).Find(&matches).Error
	}

	return matches, query.Find(&matches).Error
}

func (r *ExchangeRequestRepository) GetMatchByID(id string) (*models.ExchangeMatch, error) {
	exchangeMatch := &models.ExchangeMatch{}
	return exchangeMatch, r.db.Where("id = ?", id).First(exchangeMatch).Error
}

func (r *ExchangeRequestRepository) UpdateMatch(match *models.ExchangeMatch) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(match).Error
	// return r.db.Model(match).Updates(map[string]interface{}{"request1_accept": match.Request1Accept, "request2_accept": match.Request2Accept}).Error
}

func (r *ExchangeRequestRepository) GetActiveExchangeRequestsByBookID(id string, userID string) ([]*models.ExchangeRequest, error) {
	exchanges := []*models.ExchangeRequest{}

	// Use a WHERE EXISTS or IN approach to get ExchangeRequests where the OfferedBooks contain the specified book_id
	err := r.db.Preload("OfferedBooks.Book").
		Where("user_google_id = ?", userID).
		Where("EXISTS (SELECT 1 FROM offered_books ob WHERE ob.exchange_request_id = exchange_requests.id AND ob.book_id = ?)", id).
		Find(&exchanges).Error

	return exchanges, err
}
