package unit

import (
	"errors"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/geo"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExchangeService_Request(t *testing.T) {
	mockRepo := new(MockExchangeRequestRepository)
	service := services.NewExchangeService(mockRepo)

	userID := "user123"
	userEmail := "test@example.com"
	desiredBookID := "book456"
	userBookIDs := []string{"book123", "book789"}
	latitude, longitude := 52.2297, 21.0122
	exchangeID := "123"

	exchange := &models.ExchangeRequest{
		UserGoogleId:  userID,
		DesiredBookID: desiredBookID,
		OfferedBooks:  []models.OfferedBook{{BookID: "book123"}, {BookID: "book789"}},
		Status:        models.ExchangeRequestStatusActive,
		Latitude:      latitude,
		Longitude:     longitude,
	}

	t.Run("Create -- Success", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(nil)
		mockRepo.On("Get", mock.Anything, userID).Return(exchange, nil)

		result, err := service.Create(userID, userEmail, desiredBookID, userBookIDs, latitude, longitude)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.UserGoogleId)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Create -- Validation Error", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(models.ErrExchangeRequestNoOfferedBooksProvided)
		_, err := service.Create(userID, userEmail, desiredBookID, []string{}, latitude, longitude)
		assert.Error(t, err)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Create -- Repository Error", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(errors.New("database error"))
		mockRepo.On("Get", mock.Anything, userID).Return(exchange, nil)

		_, err := service.Create(userID, userEmail, desiredBookID, userBookIDs, latitude, longitude)

		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Create -- Repository Get Error", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything).Return(nil)
		mockRepo.On("Get", mock.Anything, userID).Return(&models.ExchangeRequest{}, errors.New("not found"))

		_, err := service.Create(userID, userEmail, desiredBookID, userBookIDs, latitude, longitude)

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Delete -- Success", func(t *testing.T) {
		mockRepo.On("GetByID", exchangeID).Return(&models.ExchangeRequest{ID: 123, Status: models.ExchangeRequestStatusActive}, nil).Once()
		mockRepo.On("Delete", exchangeID).Return(nil).Once()
		mockRepo.On("DeleteMatchesForRequest", exchangeID).Return(nil).Once()

		err := service.Delete(exchangeID)

		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Delete -- Exchange Not Found", func(t *testing.T) {
		mockRepo.On("GetByID", exchangeID).Return(&models.ExchangeRequest{}, errors.New("not found")).Once()

		err := service.Delete(exchangeID)

		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())

		mockRepo.AssertCalled(t, "GetByID", exchangeID)
		mockRepo.AssertNotCalled(t, "Delete")
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Delete -- Exchange Status Completed", func(t *testing.T) {
		mockRepo.On("GetByID", exchangeID).Return(&models.ExchangeRequest{ID: 123, Status: models.ExchangeRequestStatusCompleted}, nil).Once()

		err := service.Delete(exchangeID)

		assert.Error(t, err)
		assert.Equal(t, models.ErrExchangeRequestCompleted, err)

		mockRepo.AssertCalled(t, "GetByID", exchangeID)
		mockRepo.AssertNotCalled(t, "Delete")
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Delete -- Repository Delete Failure", func(t *testing.T) {
		mockRepo.On("GetByID", exchangeID).Return(&models.ExchangeRequest{ID: 123, Status: models.ExchangeRequestStatusActive}, nil).Once()
		mockRepo.On("Delete", exchangeID).Return(errors.New("delete failed")).Once()

		err := service.Delete(exchangeID)

		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())

		mockRepo.AssertCalled(t, "GetByID", exchangeID)
		mockRepo.AssertCalled(t, "Delete", exchangeID)
		mockRepo.AssertNotCalled(t, "DeleteMatchesForRequest")
		mockRepo.ClearExpectedCalls()
	})
}

func TestExchangeService_Matches(t *testing.T) {
	mockRepo := new(MockExchangeRequestRepository)
	service := services.NewExchangeService(mockRepo)

	requestID := "1"
	userID := "user123"
	request, otherRequests, expectedMatch := seedDataRequestsMatches(t)
	matches := seedDataFilterMatches(t)
	distanceThreshold := 10.0

	t.Run("Successful Match Creation", func(t *testing.T) {
		mockRepo.On("CreateMatch", expectedMatch).Return(nil)
		match, err := service.CreateMatch(request, otherRequests[0])

		assert.NoError(t, err)
		assert.Equal(t, expectedMatch, match)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Successful Distance Filtering", func(t *testing.T) {
		mockRepo.On("GetAllMatches", requestID).Return(matches, nil)
		expectedFilteredMatches := []*models.ExchangeMatch{matches[0], matches[1]}
		filteredMatches, err := service.GetMatchesDistanceFiltered(requestID, distanceThreshold)

		assert.NoError(t, err)
		assert.Equal(t, expectedFilteredMatches, filteredMatches)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("No Matches Found", func(t *testing.T) {
		mockRepo.On("GetAllMatches", requestID).Return([]*models.ExchangeMatch{}, nil)

		filteredMatches, err := service.GetMatchesDistanceFiltered(requestID, distanceThreshold)

		assert.NoError(t, err)
		assert.Empty(t, filteredMatches)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("All Matches Exceed Threshold", func(t *testing.T) {
		highDistanceMatches := []*models.ExchangeMatch{
			{ExchangeRequestID: 1, MatchedExchangeRequestID: 6, Distance: 20.0},
			{ExchangeRequestID: 1, MatchedExchangeRequestID: 7, Distance: 30.0},
		}
		mockRepo.On("GetAllMatches", requestID).Return(highDistanceMatches, nil)

		filteredMatches, err := service.GetMatchesDistanceFiltered(requestID, distanceThreshold)

		assert.NoError(t, err)
		assert.Empty(t, filteredMatches)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Repository Returns Error", func(t *testing.T) {
		mockRepo.On("GetAllMatches", requestID).Return([]*models.ExchangeMatch{}, errors.New("database error"))

		filteredMatches, err := service.GetMatchesDistanceFiltered(requestID, distanceThreshold)

		assert.Error(t, err)
		assert.Nil(t, filteredMatches)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("Successful Match Finding", func(t *testing.T) {
		mockRepo.On("Get", requestID, userID).Return(request, nil)
		mockRepo.On("FindMatchingRequests", requestID, userID, request.DesiredBookID, getOfferedBookIDs(request.OfferedBooks)).Return(otherRequests, nil)
		mockRepo.On("Update", request).Return(nil)
		mockRepo.On("CreateMatch", expectedMatch).Return(nil, nil)

		matchedRequests, err := service.FindMatchingRequests(requestID, userID)

		assert.NoError(t, err)
		assert.Equal(t, otherRequests, matchedRequests)
		mockRepo.AssertExpectations(t)
		mockRepo.ClearExpectedCalls()
	})
}

func seedDataRequestsMatches(t *testing.T) (*models.ExchangeRequest, []*models.ExchangeRequest, *models.ExchangeMatch) {
	t.Helper()
	userID := "user123"
	request := &models.ExchangeRequest{
		ID:            1,
		UserGoogleId:  userID,
		DesiredBookID: "bookA",
		OfferedBooks:  []models.OfferedBook{{BookID: "bookB"}},
		Status:        models.ExchangeRequestStatusActive,
		Latitude:      0,
		Longitude:     0,
	}
	otherRequests := []*models.ExchangeRequest{
		{
			ID:            2,
			UserGoogleId:  "user456",
			DesiredBookID: "bookB",
			OfferedBooks:  []models.OfferedBook{{BookID: "bookA"}},
			Status:        models.ExchangeRequestStatusActive,
			Latitude:      0,
		},
	}

	expectedDistance := geo.HaversineDistance(
		geo.Cord{Lat: request.Latitude, Lon: request.Longitude},
		geo.Cord{Lat: otherRequests[0].Latitude, Lon: otherRequests[0].Longitude},
		geo.Km,
	)

	expectedMatch := &models.ExchangeMatch{
		ExchangeRequestID:        1,
		MatchedExchangeRequestID: 2,
		Status:                   models.MatchStatusPending,
		Request1Decision:         models.MatchDecisionPending,
		Request2Decision:         models.MatchDecisionPending,
		Distance:                 expectedDistance,
	}
	return request, otherRequests, expectedMatch
}

func TestExchangeService_AcceptMatch(t *testing.T) {
	mockRepo := new(MockExchangeRequestRepository)
	service := services.NewExchangeService(mockRepo)

	pendingMatch, oneAcceptedMatch, acceptedMatch := seedDataMatches(t)

	testCases := []struct {
		name          string
		matchID       string
		requestID     string
		mockMatch     *models.ExchangeMatch
		mockError     error
		expectedMatch *models.ExchangeMatch
		expectError   bool
	}{
		{
			name:          "First user accepts - match still pending",
			matchID:       "match123",
			requestID:     "1",
			mockMatch:     pendingMatch,
			expectedMatch: oneAcceptedMatch,
			expectError:   false,
		},
		{
			name:          "Second user accepts - match finalized",
			matchID:       "match123",
			requestID:     "2",
			mockMatch:     oneAcceptedMatch,
			expectedMatch: acceptedMatch,
			expectError:   false,
		},
		{
			name:        "Match not found",
			matchID:     "invalid123",
			requestID:   "1",
			mockMatch:   nil,
			mockError:   errors.New("match not found"),
			expectError: true,
		},
		{
			name:        "Repository update fails",
			matchID:     "match123",
			requestID:   "1",
			mockMatch:   pendingMatch,
			mockError:   errors.New("database update error"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer mockRepo.ClearExpectedCalls()
			mockRepo.On("GetMatch", tc.matchID, tc.requestID).Return(tc.mockMatch, tc.mockError)

			if tc.mockMatch != nil && tc.mockError == nil {
				mockRepo.On("UpdateMatch", mock.Anything).Return(tc.mockError)
			}

			match, err := service.AcceptMatch(tc.matchID, tc.requestID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, match)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMatch, match)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestExchangeService_DeclineMatch(t *testing.T) {
	mockRepo := new(MockExchangeRequestRepository)
	service := services.NewExchangeService(mockRepo)

	matchID := "123"
	requestID := "1"

	_, _, match := seedDataMatches(t)

	declinedMatch := *match
	declinedMatch.Status = models.MatchStatusDeclined
	declinedMatch.Request1Decision = models.MatchDecisionDeclined

	testCases := []struct {
		name          string
		mockSetup     func()
		expectedMatch *models.ExchangeMatch
		expectedErr   bool
	}{
		{
			name: "Successful Decline",
			mockSetup: func() {
				mockRepo.On("GetMatch", matchID, requestID).Return(match, nil)
				mockRepo.On("UpdateMatch", &declinedMatch).Return(nil)
			},
			expectedMatch: &declinedMatch,
			expectedErr:   false,
		},
		{
			name: "Match Already Declined",
			mockSetup: func() {
				alreadyDeclinedMatch := *match
				alreadyDeclinedMatch.Request1Decision = models.MatchDecisionDeclined
				alreadyDeclinedMatch.Status = models.MatchStatusDeclined

				mockRepo.On("GetMatch", matchID, requestID).Return(&alreadyDeclinedMatch, nil)
				mockRepo.On("UpdateMatch", &alreadyDeclinedMatch).Return(nil)
			},
			expectedMatch: &declinedMatch,
			expectedErr:   false,
		},
		{
			name: "Error Retrieving Match",
			mockSetup: func() {
				mockRepo.On("GetMatch", matchID, requestID).Return(&models.ExchangeMatch{}, errors.New("match not found"))
			},
			expectedMatch: nil,
			expectedErr:   true,
		},
		{
			name: "Error Updating Match",
			mockSetup: func() {
				mockRepo.On("GetMatch", matchID, requestID).Return(match, nil)
				mockRepo.On("UpdateMatch", &declinedMatch).Return(errors.New("update failed"))
			},
			expectedMatch: nil,
			expectedErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer mockRepo.ClearExpectedCalls()
			tc.mockSetup()

			result, err := service.DeclineMatch(matchID, requestID)

			if tc.expectedErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedMatch, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func seedDataMatches(t *testing.T) (*models.ExchangeMatch, *models.ExchangeMatch, *models.ExchangeMatch) {
	t.Helper()
	pendingMatch := &models.ExchangeMatch{
		ID:                       123,
		ExchangeRequestID:        1,
		MatchedExchangeRequestID: 2,
		Status:                   models.MatchStatusPending,
		Request1Decision:         models.MatchDecisionPending,
		Request2Decision:         models.MatchDecisionPending,
	}

	oneAcceptedMatch := *pendingMatch
	oneAcceptedMatch.Request1Decision = models.MatchDecisionAccepted

	acceptedMatch := *pendingMatch
	acceptedMatch.Request1Decision = models.MatchDecisionAccepted
	acceptedMatch.Request2Decision = models.MatchDecisionAccepted
	acceptedMatch.Status = models.MatchStatusAccepted
	acceptedMatch.Request.Status = models.ExchangeRequestStatusCompleted
	acceptedMatch.MatchedExchangeRequest.Status = models.ExchangeRequestStatusCompleted

	return pendingMatch, &oneAcceptedMatch, &acceptedMatch
}

func seedDataFilterMatches(t *testing.T) []*models.ExchangeMatch {
	t.Helper()
	return []*models.ExchangeMatch{
		{ExchangeRequestID: 1, MatchedExchangeRequestID: 2, Distance: 5.0},
		{ExchangeRequestID: 1, MatchedExchangeRequestID: 3, Distance: 9.9},
		{ExchangeRequestID: 1, MatchedExchangeRequestID: 4, Distance: 10.1},
		{ExchangeRequestID: 1, MatchedExchangeRequestID: 5, Distance: 15.0},
	}
}

func getOfferedBookIDs(offeredBooks []models.OfferedBook) []string {
	ids := make([]string, len(offeredBooks))
	for i, book := range offeredBooks {
		ids[i] = book.BookID
	}
	return ids
}
