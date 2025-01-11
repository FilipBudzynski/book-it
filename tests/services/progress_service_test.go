package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
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

func TestCreateProgress(t *testing.T) {
	mockRepo, service := setupMockRepoAndService(t)
	t.Run("Valid Progress Creation", func(t *testing.T) {
		startDate := "2025-01-01"
		endDate := "2025-01-10"
		bookId := uint(1)
		totalPages := 100
		bookTitle := "Test Book"

		// Mock expectations
		mockRepo.On("Create", mock.Anything).Return(nil)

		progress, err := service.Create(bookId, totalPages, bookTitle, startDate, endDate)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, bookId, progress.UserBookID)
		assert.Equal(t, totalPages, progress.TotalPages)
		assert.Equal(t, 10, len(progress.DailyProgress))
		assert.Equal(t, 10, progress.DailyTargetPages)
		mockRepo.AssertCalled(t, "Create", mock.Anything)
	})

	t.Run("Invalid Start Date", func(t *testing.T) {
		_, err := service.Create(1, 100, "Test Book", "invalid-date", "2025-01-10")
		assert.Error(t, err)
	})

	t.Run("End Date Before Start Date", func(t *testing.T) {
		_, err := service.Create(1, 100, "Test Book", "2025-01-10", "2025-01-01")
		assert.Equal(t, models.ErrProgressInvalidEndDate, err)
	})
}

func setupMockRepoAndService(t *testing.T) (*MockProgressRepository, handlers.ProgressService) {
	t.Helper()
	mockRepo := new(MockProgressRepository)
	service := services.NewProgressService(mockRepo)
	return mockRepo, service
}

func TestUpdateTargetPages(t *testing.T) {
	t.Run("Valid Target Pages Update", func(t *testing.T) {
		mockRepo, service := setupMockRepoAndService(t)
		logDate := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC)

		mockProgress := setupMockProgress(t, 20, 100, 5, logDate)
		mockProgress.EndDate = endDate

		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)

		err := service.UpdateTargetPages(1)

		assert.NoError(t, err)
		assert.Equal(t, 20, mockProgress.DailyTargetPages)
		for _, log := range mockProgress.DailyProgress {
			assert.Equal(t, 20, log.TargetPages)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("Valid Update with Positive Days Left", func(t *testing.T) {
		logDate := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
		mockProgress := setupMockProgress(t, 50, 100, 6, logDate)

		mockRepo, service := setupMockRepoAndService(t)
		mockProgress.EndDate = endDate

		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)

		err := service.UpdateTargetPages(1)

		assert.NoError(t, err)
		assert.Equal(t, 10, mockProgress.DailyTargetPages)
		for _, log := range mockProgress.DailyProgress {
			assert.Equal(t, 10, log.TargetPages)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update in the middle of progress", func(t *testing.T) {
		startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		logDate := time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
		mockProgress := setupMockProgress(t, 30, 100, 10, startDate)

		// Set completed logs
		for i := 0; i < 5; i++ {
			mockProgress.DailyProgress[i].TargetPages = 10
			mockProgress.DailyProgress[i].PagesRead = 5
		}

		mockRepo, service := setupMockRepoAndService(t)
		mockProgress.EndDate = endDate

		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)

		err := service.UpdateTargetPages(1)

		assert.NoError(t, err)
		assert.Equal(t, 14, mockProgress.DailyTargetPages)
		for _, log := range mockProgress.DailyProgress {
			if log.Date.Before(logDate) {
				assert.Equal(t, 10, log.TargetPages)
			} else {
				assert.Equal(t, 14, log.TargetPages)
			}
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("Last Day - book not finished", func(t *testing.T) {
		logDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
		mockProgress := setupMockProgress(t, 50, 55, 1, logDate)

		mockRepo, service := setupMockRepoAndService(t)
		mockProgress.EndDate = endDate

		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)

		_ = service.UpdateTargetPages(1)

		assert.Equal(t, 5, mockProgress.DailyTargetPages)
		mockRepo.AssertExpectations(t)
	})
}

func setupMockProgress(t *testing.T, currentPage, totalPages, days int, startDate time.Time) *models.ReadingProgress {
	t.Helper()
	progress := &models.ReadingProgress{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		DailyProgress: make([]models.DailyProgressLog, days),
	}
	for i := 0; i < days; i++ {
		progress.DailyProgress[i] = models.DailyProgressLog{
			Date: startDate.AddDate(0, 0, i),
		}
	}
	return progress
}

func TestUpdateLogPagesRead(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := services.NewProgressService(mockRepo)

	t.Run("Valid Log Update", func(t *testing.T) {
		logId := "1"
		pagesRead := 10
		mockTime := time.Now()
		mockLog := &models.DailyProgressLog{
			ReadingProgressID: 1,
			Date:              mockTime,
		}

		mockProgress := setupMockProgress(t, 50, 100, 10, mockTime)
		mockProgress.EndDate = time.Now().Add(10 * time.Hour * 24)

		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)
		mockRepo.On("GetLogById", logId).Return(mockLog, nil)
		mockRepo.On("UpdateLog", mockLog).Return(nil)

		_, err := service.UpdateLog(logId, pagesRead, "")

		assert.NoError(t, err)
		assert.Equal(t, pagesRead, mockLog.PagesRead)
		mockRepo.AssertCalled(t, "UpdateLog", mockLog)
	})

	t.Run("Pages read not changed", func(t *testing.T) {
		logId := "1"
		pagesRead := 0
		mockLog := &models.DailyProgressLog{
			PagesRead:         0,
			ReadingProgressID: 1,
			Date:              time.Now(),
		}

		mockProgress := setupMockProgress(t, 50, 100, 10, time.Now())
		mockRepo.On("GetById", "1").Return(mockProgress, nil)
		mockRepo.On("Update", mockProgress).Return(nil)
		mockRepo.On("GetLogById", logId).Return(mockLog, nil)
		mockRepo.On("UpdateLog", mockLog).Return(nil)

		_, err := service.UpdateLog(logId, pagesRead, "")

		assert.NoError(t, err)
		assert.Equal(t, 0, mockLog.PagesRead)
		mockRepo.AssertNotCalled(t, "UpdateLog", mockLog)
	})

	t.Run("Negative PagesRead Input", func(t *testing.T) {
		_, err := service.UpdateLog("1", -1, "")
		assert.Error(t, err)
	})
}

func TestGetProgress(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := services.NewProgressService(mockRepo)

	mockProgress := &models.ReadingProgress{
		UserBookID: 1,
	}

	// Mock expectations
	mockRepo.On("GetById", "1").Return(mockProgress, nil)

	progress, err := service.Get("1")

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, uint(1), progress.UserBookID)
	mockRepo.AssertCalled(t, "GetById", "1")
}

func TestGetProgress_NotFound(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := services.NewProgressService(mockRepo)

	// Mock expectations
	mockRepo.On("GetById", "1").Return(nil, errors.New("not found"))

	progress, err := service.Get("1")

	// Assertions
	assert.Nil(t, progress)
	assert.Error(t, err)
	mockRepo.AssertCalled(t, "GetById", "1")
}
