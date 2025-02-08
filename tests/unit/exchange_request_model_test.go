package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestExchangeRequest_GetMatchStatus(t *testing.T) {
	tests := []struct {
		name           string
		exchangeReq    models.ExchangeRequest
		otherRequestID uint
		expectedStatus models.MatchStatus
	}{
		{
			name: "Match found - MatchedExchangeRequestID",
			exchangeReq: models.ExchangeRequest{
				Matches: []models.ExchangeMatch{
					{ExchangeRequestID: 1, MatchedExchangeRequestID: 2, Status: models.MatchStatusAccepted},
				},
			},
			otherRequestID: 2,
			expectedStatus: models.MatchStatusAccepted,
		},
		{
			name: "Match found - ExchangeRequestID",
			exchangeReq: models.ExchangeRequest{
				Matches: []models.ExchangeMatch{
					{ExchangeRequestID: 3, MatchedExchangeRequestID: 4, Status: models.MatchStatusDeclined},
				},
			},
			otherRequestID: 3,
			expectedStatus: models.MatchStatusDeclined,
		},
		{
			name: "No matching request",
			exchangeReq: models.ExchangeRequest{
				Matches: []models.ExchangeMatch{
					{ExchangeRequestID: 5, MatchedExchangeRequestID: 6, Status: models.MatchStatusAccepted},
				},
			},
			otherRequestID: 7,
			expectedStatus: models.MatchStatusPending,
		},
		{
			name: "Empty matches list",
			exchangeReq: models.ExchangeRequest{
				Matches: []models.ExchangeMatch{},
			},
			otherRequestID: 1,
			expectedStatus: models.MatchStatusPending,
		},
		{
			name:           "Nil matches list",
			exchangeReq:    models.ExchangeRequest{},
			otherRequestID: 1,
			expectedStatus: models.MatchStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.exchangeReq.GetMatchStatus(tt.otherRequestID)
			assert.Equal(t, tt.expectedStatus, status)
		})
	}
}

func TestExchangeRequest_Validate(t *testing.T) {
	tests := []struct {
		name          string
		exchangeReq   models.ExchangeRequest
		expectedError error
	}{
		{
			name: "Valid exchange request",
			exchangeReq: models.ExchangeRequest{
				DesiredBookID: "123",
				OfferedBooks:  []models.OfferedBook{{BookID: "123"}, {BookID: "456"}},
				Latitude:      40.7128,
				Longitude:     -74.0060,
			},
			expectedError: nil,
		},
		{
			name: "Missing desired book",
			exchangeReq: models.ExchangeRequest{
				OfferedBooks: []models.OfferedBook{{ID: 1}, {ID: 2}},
				Latitude:     40.7128,
				Longitude:    -74.0060,
			},
			expectedError: models.ErrExchangeRequestNoDesiredBookProvided,
		},
		{
			name: "No offered books",
			exchangeReq: models.ExchangeRequest{
				DesiredBookID: "123",
				Latitude:      40.7128,
				Longitude:     -74.0060,
			},
			expectedError: models.ErrExchangeRequestNoOfferedBooksProvided,
		},
		{
			name: "Latitude out of range",
			exchangeReq: models.ExchangeRequest{
				DesiredBookID: "123",
				OfferedBooks:  []models.OfferedBook{{ID: 1}},
				Latitude:      95.0,
				Longitude:     -74.0060,
			},
			expectedError: models.ErrExchangeRequestLatitudeOutOfRange,
		},
		{
			name: "Longitude out of range",
			exchangeReq: models.ExchangeRequest{
				DesiredBookID: "123",
				OfferedBooks:  []models.OfferedBook{{ID: 1}},
				Latitude:      40.7128,
				Longitude:     200.0,
			},
			expectedError: models.ErrExchangeRequestLongitudeOutOfRange,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.exchangeReq.Validate()
			if tt.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}

func TestExchangeMatch(t *testing.T) {
	tests := []struct {
		name          string
		exchangeMatch models.ExchangeMatch
		requestID     uint
		expectedMatchedRequest *models.ExchangeRequest
		expectedAccepted       bool
		expectedDecision      models.MatchDecision
	}{
		{
			name: "Request1 accepted, matching request and decision for Request1",
			exchangeMatch: models.ExchangeMatch{
				ExchangeRequestID:  1,
				MatchedExchangeRequestID: 2,
				Request1Decision:   models.MatchDecisionAccepted,
				Request2Decision:   models.MatchDecisionPending,
				Request:             models.ExchangeRequest{ID: 1},
				MatchedExchangeRequest: models.ExchangeRequest{ID: 2},
			},
			requestID:            1,
			expectedMatchedRequest: &models.ExchangeRequest{ID: 2},
			expectedAccepted:      true,
			expectedDecision:      models.MatchDecisionAccepted,
		},
		{
			name: "Request2 accepted, matching request and decision for Request2",
			exchangeMatch: models.ExchangeMatch{
				ExchangeRequestID:  1,
				MatchedExchangeRequestID: 2,
				Request1Decision:   models.MatchDecisionPending,
				Request2Decision:   models.MatchDecisionAccepted,
				Request:             models.ExchangeRequest{ID: 1},
				MatchedExchangeRequest: models.ExchangeRequest{ID: 2},
			},
			requestID:            2,
			expectedMatchedRequest: &models.ExchangeRequest{ID: 1},
			expectedAccepted:      true,
			expectedDecision:      models.MatchDecisionAccepted,
		},
		{
			name: "Neither request accepted, decision for Request1 is pending",
			exchangeMatch: models.ExchangeMatch{
				ExchangeRequestID:  1,
				MatchedExchangeRequestID: 2,
				Request1Decision:   models.MatchDecisionPending,
				Request2Decision:   models.MatchDecisionDeclined,
				Request:             models.ExchangeRequest{ID: 1},
				MatchedExchangeRequest: models.ExchangeRequest{ID: 2},
			},
			requestID:            1,
			expectedMatchedRequest: &models.ExchangeRequest{ID: 2},
			expectedAccepted:      false,
			expectedDecision:      models.MatchDecisionPending,
		},
		{
			name: "Request1 accepted, decision for Request2 is declined",
			exchangeMatch: models.ExchangeMatch{
				ExchangeRequestID:  1,
				MatchedExchangeRequestID: 2,
				Request1Decision:   models.MatchDecisionAccepted,
				Request2Decision:   models.MatchDecisionDeclined,
				Request:             models.ExchangeRequest{ID: 1},
				MatchedExchangeRequest: models.ExchangeRequest{ID: 2},
			},
			requestID:            2,
			expectedMatchedRequest: &models.ExchangeRequest{ID: 1},
			expectedAccepted:      false,
			expectedDecision:      models.MatchDecisionDeclined,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("MatchedRequest", func(t *testing.T) {
				result := tt.exchangeMatch.MatchedRequest(tt.requestID)
				assert.Equal(t, tt.expectedMatchedRequest, result)
			})

			t.Run("IsAccepted", func(t *testing.T) {
				result := tt.exchangeMatch.IsAccepted(tt.requestID)
				assert.Equal(t, tt.expectedAccepted, result)
			})

			t.Run("GetDecision", func(t *testing.T) {
				result := tt.exchangeMatch.GetDecision(tt.requestID)
				assert.Equal(t, tt.expectedDecision, result)
			})
		})
	}
}
