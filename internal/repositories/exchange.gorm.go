package repositories

import (
	"fmt"
	"strconv"

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

func (r *ExchangeRequestRepository) FindMatchingRequests(
	requestId,
	userId,
	desiredBookId string,
	offeredBooks []string,
	existingMatches []models.ExchangeMatch,
) ([]models.ExchangeMatch, error) {
	var existingMatchIDs []uint
	for _, match := range existingMatches {
		existingMatchIDs = append(existingMatchIDs, match.MatchRequestID)
	}

	query := r.db.Model(&models.ExchangeRequest{}).
		Not("id = ?", requestId).                                                                        // Exclude user's own request
		Not("user_google_id = ?", userId).                                                               // Exclude user's own requests
		Where("desired_book_id IN ?", offeredBooks).                                                     // Their desired book is in your offered books
		Where("id IN (SELECT exchange_request_id FROM offered_books WHERE book_id = ?)", desiredBookId). // Your desired book is in their offered books
		Where("status = ?", models.ExchangeRequestStatusPending)                                         // Only match pending requests

	if len(existingMatchIDs) > 0 {
		query = query.Not("id IN ?", existingMatchIDs)
	}

	var matchIDs []uint
	err := query.Pluck("id", &matchIDs).Error
	if err != nil {
		return nil, err
	}

	fmt.Printf("MATCHES LEN REPO: %d\n", len(matchIDs))

	requestIdInt, _ := strconv.Atoi(requestId)
	requestIdUint := uint(requestIdInt)
	matchesFound := make([]models.ExchangeMatch, len(matchIDs))
	for i, id := range matchIDs {
		match := models.ExchangeMatch{
			ExchangeRequestID: requestIdUint,
			MatchRequestID:    id,
		}
		matchesFound[i] = match
	}
	return matchesFound, nil
}

// var matches []*models.ExchangeRequest
// err := r.db.Preload("DesiredBook").
// 	Preload("OfferedBooks.Book").
// 	Not("id = ?", requestId).                                                                        // Exclude user's own request
// 	Not("user_google_id = ?", userId).                                                               // Exclude user's own requests
// 	Where("desired_book_id IN ?", offeredBooks).                                                     // Their desired book is in your offered books
// 	Where("id IN (SELECT exchange_request_id FROM offered_books WHERE book_id = ?)", desiredBookId). // Your desired book is in their offered books
// 	Where("status = ?", models.ExchangeRequestStatusPending).                                        // Only match pending requests
// 	Find(&matches).Error
// if err != nil {
// 	return nil, err
// }
