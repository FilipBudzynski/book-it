package unit

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockProgressRepository struct {
	mock.Mock
}

func (m *MockProgressRepository) Create(progress models.ReadingProgress) error {
	args := m.Called(progress)
	return args.Error(0)
}

func (m *MockProgressRepository) GetById(id string) (*models.ReadingProgress, error) {
	args := m.Called(id)
	return args.Get(0).(*models.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) GetByUserBookId(userBookId string) (*models.ReadingProgress, error) {
	args := m.Called(userBookId)
	return args.Get(0).(*models.ReadingProgress), args.Error(1)
}

func (m *MockProgressRepository) Update(progress *models.ReadingProgress) error {
	args := m.Called(progress)
	return args.Error(0)
}

func (m *MockProgressRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProgressRepository) GetLogById(id string) (*models.DailyProgressLog, error) {
	args := m.Called(id)
	return args.Get(0).(*models.DailyProgressLog), args.Error(1)
}

func (m *MockProgressRepository) UpdateLog(log *models.DailyProgressLog) error {
	args := m.Called(log)
	return args.Error(0)
}

type MockProgressService struct {
	mock.Mock
}

func (m *MockProgressService) UpdateTargetPages(progress *models.ReadingProgress, logID uint) (*models.ReadingProgress, error) {
	args := m.Called(progress, logID)
	return args.Get(0).(*models.ReadingProgress), args.Error(1)
}
