package unit

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockExchangeRequestRepository struct {
	mock.Mock
}

func (m *MockExchangeRequestRepository) Create(exchange *models.ExchangeRequest) error {
	args := m.Called(exchange)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) Get(id, userId string) (*models.ExchangeRequest, error) {
	args := m.Called(id, userId)
	return args.Get(0).(*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) GetByID(id string) (*models.ExchangeRequest, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) GetAll(userId string) ([]*models.ExchangeRequest, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) GetAllWithStatus(userId string, status models.ExchangeRequestStatus) ([]*models.ExchangeRequest, error) {
	args := m.Called(userId, status)
	return args.Get(0).([]*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) DeleteMatchesForRequest(requestId string) error {
	args := m.Called(requestId)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) Update(exchange *models.ExchangeRequest) error {
	args := m.Called(exchange)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) FindMatchingRequests(userId, requestId, desiredBookId string, offeredBooks []string) ([]*models.ExchangeRequest, error) {
	args := m.Called(userId, requestId, desiredBookId, offeredBooks)
	return args.Get(0).([]*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) CreateMatch(match *models.ExchangeMatch) error {
	args := m.Called(match)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) GetMatch(requestId, otherRequestId string) (*models.ExchangeMatch, error) {
	args := m.Called(requestId, otherRequestId)
	return args.Get(0).(*models.ExchangeMatch), args.Error(1)
}

func (m *MockExchangeRequestRepository) UpdateMatch(match *models.ExchangeMatch) error {
	args := m.Called(match)
	return args.Error(0)
}

func (m *MockExchangeRequestRepository) GetAllMatches(requestId string) ([]*models.ExchangeMatch, error) {
	args := m.Called(requestId)
	return args.Get(0).([]*models.ExchangeMatch), args.Error(1)
}

func (m *MockExchangeRequestRepository) GetMatchByID(id string) (*models.ExchangeMatch, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ExchangeMatch), args.Error(1)
}

func (m *MockExchangeRequestRepository) GetActiveExchangeRequestsByBookID(id string, userID string) ([]*models.ExchangeRequest, error) {
	args := m.Called(id, userID)
	return args.Get(0).([]*models.ExchangeRequest), args.Error(1)
}

func (m *MockExchangeRequestRepository) ClearExpectedCalls() {
	m.Mock.ExpectedCalls = nil
}
