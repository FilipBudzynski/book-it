package unit

import (
	"errors"
	"math"
	"testing"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProgressService(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	service := services.NewProgressService(mockRepo)

	t.Run("Create simple", func(t *testing.T) {
		tests := []struct {
			name        string
			bookId      uint
			totalPages  int
			bookTitle   string
			startDate   string
			endDate     string
			mockReturn  error
			expectError bool
		}{
			{
				name:       "Valid input",
				bookId:     1,
				totalPages: 300,
				bookTitle:  "Test Book",
				startDate:  "2024-02-01",
				endDate:    "2024-02-10",
				mockReturn: nil,
			},
			{
				name:        "Invalid date range (end before start)",
				bookId:      2,
				totalPages:  100,
				bookTitle:   "Test Book 2",
				startDate:   "2024-02-10",
				endDate:     "2024-02-01",
				mockReturn:  models.ErrProgressInvalidEndDate,
				expectError: true,
			},
			{
				name:        "Exceed max daily logs",
				bookId:      3,
				totalPages:  500,
				bookTitle:   "Long Read",
				startDate:   "2024-02-01",
				endDate:     "2025-02-01",
				mockReturn:  models.ErrProgressMaxLogsExceeded,
				expectError: true,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				mockRepo.On("Create", mock.Anything).Return(tc.mockReturn)

				progress, err := service.Create(tc.bookId, tc.totalPages, tc.bookTitle, tc.startDate, tc.endDate)

				if tc.expectError {
					assert.Error(t, err)
					assert.Equal(t, tc.mockReturn, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.bookId, progress.UserBookID)
					assert.Equal(t, tc.totalPages, progress.TotalPages)
					assert.Equal(t, tc.bookTitle, progress.BookTitle)
				}

				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("Create with logs", func(t *testing.T) {
		startDate := "2024-02-01"
		endDate := "2024-02-05"
		bookId := uint(1)
		totalPages := 100
		bookTitle := "Test Book"

		expectedStartDate, _ := time.Parse(time.DateOnly, startDate)
		expectedEndDate, _ := time.Parse(time.DateOnly, endDate)
		days := int(expectedEndDate.Sub(expectedStartDate).Hours()/24) + 1

		mockRepo.On("Create", mock.Anything).Return(nil)

		progress, err := service.Create(bookId, totalPages, bookTitle, startDate, endDate)

		assert.NoError(t, err)
		assert.Equal(t, bookId, progress.UserBookID)
		assert.Equal(t, totalPages, progress.TotalPages)
		assert.Equal(t, bookTitle, progress.BookTitle)
		assert.Equal(t, expectedStartDate, progress.StartDate)
		assert.Equal(t, expectedEndDate, progress.EndDate)

		assert.Len(t, progress.DailyProgress, days, "Should generate correct number of logs")

		pagesLeft := totalPages
		for i, log := range progress.DailyProgress {
			expectedDate := expectedStartDate.AddDate(0, 0, i)
			assert.Equal(t, expectedDate, log.Date, "Log date should be sequential")
			assert.Equal(t, bookId, log.UserBookID, "Log should be linked to correct book ID")

			expectedPages := (pagesLeft + (days - i) - 1) / (days - i)
			assert.Equal(t, expectedPages, log.TargetPages, "Target pages should be correctly distributed")
			pagesLeft -= expectedPages
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Invalid Cases", func(t *testing.T) {
		tests := []struct {
			name        string
			bookId      uint
			totalPages  int
			bookTitle   string
			startDate   string
			endDate     string
			mockReturn  error
			expectError bool
			expectedErr error
		}{
			{
				name:        "End date before start date",
				bookId:      1,
				totalPages:  100,
				bookTitle:   "Reverse Dates",
				startDate:   "2024-02-10",
				endDate:     "2024-02-01", // Invalid: end before start
				expectError: true,
				expectedErr: models.ErrProgressInvalidEndDate,
			},
			{
				name:        "Exceed max daily logs",
				bookId:      2,
				totalPages:  500,
				bookTitle:   "Too Many Logs",
				startDate:   "2024-02-01",
				endDate:     "2025-02-01", // Exceeds max days allowed
				expectError: true,
				expectedErr: models.ErrProgressMaxLogsExceeded,
			},
			{
				name:        "Zero pages not allowed",
				bookId:      3,
				totalPages:  0,
				bookTitle:   "Zero Pages",
				startDate:   "2024-02-01",
				endDate:     "2024-02-05",
				expectError: true,
				expectedErr: models.ErrProgressInvalidTotalPages,
			},
			{
				name:        "Negative pages not allowed",
				bookId:      4,
				totalPages:  -50,
				bookTitle:   "Negative Pages",
				startDate:   "2024-02-01",
				endDate:     "2024-02-10",
				expectError: true,
				expectedErr: models.ErrProgressInvalidTotalPages,
			},
			{
				name:        "Invalid date format",
				bookId:      5,
				totalPages:  100,
				bookTitle:   "Bad Date Format",
				startDate:   "Feb 01, 2024",
				endDate:     "Feb 05, 2024",
				expectError: true,
				expectedErr: &time.ParseError{Layout: "2006-01-02", Value: "Feb 01, 2024", LayoutElem: "2006", ValueElem: "Feb 01, 2024", Message: ""},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				mockRepo.On("Create", mock.Anything).Return(tc.mockReturn)

				progress, err := service.Create(tc.bookId, tc.totalPages, tc.bookTitle, tc.startDate, tc.endDate)

				assert.Error(t, err, "Expected an error but got none")
				assert.Equal(t, tc.expectedErr, err, "Unexpected error message")
				assert.Equal(t, models.ReadingProgress{}, progress, "Expected empty progress on failure")

				mockRepo.AssertNotCalled(t, "Create")
			})
		}
	})
}

func TestProgressService_CalculateTargetPages(t *testing.T) {
	tests := []struct {
		name      string
		pagesLeft int
		daysLeft  int
		expected  int
	}{
		{"Zero Days Left", 100, 0, 100},
		{"Even Distribution", 100, 10, 10},
		{"Rounding Up", 101, 10, 11},
		{"No Pages Left", 0, 10, 0},
		{"Negative Pages Left", -5, 10, -1},
		{"Negative Days Left", 100, -10, -1},
		{"One Day Left", 50, 1, 50},
		{"One Page Left, One Day Left", 1, 1, 1},
		{"One Page Left, Multiple Days", 1, 10, 1},
		{"Large Pages and Days", 1000000, 1000, 1000},
		{"Large Pages, Few Days", 1000000, 10, 100000},
		{"Few Pages, Large Days", 10, 1000, 1},
		{"Zero Pages and Zero Days", 0, 0, 0},
		{"Max Int Pages, One Day", math.MaxInt64, 1, math.MaxInt64},
		{"Very Uneven Distribution", 999, 500, 2},
		{"Minimal Days, Large Pages", 1000, 2, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.CalculateTargetPages(tt.pagesLeft, tt.daysLeft)
			assert.Equal(t, tt.expected, result, "Expected %d but got %d", tt.expected, result)
		})
	}
}

func TestProgressService_UpdateTargetPages(t *testing.T) {
	today := utils.TodaysDate()
	yesterday := today.AddDate(0, 0, -1)
	tomorrow := today.AddDate(0, 0, 1)

	tests := []struct {
		name     string
		progress *models.ReadingProgress
		logID    uint
		expected []int
	}{
		{
			name: "Default case, even split between the days",
			progress: &models.ReadingProgress{
				TotalPages: 100,
				EndDate:    today.AddDate(0, 0, 4),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: today, PagesRead: 0},
					{ID: 2, Date: tomorrow, PagesRead: 0},
				},
			},
			logID:    1,
			expected: []int{20, 20},
		},
		{
			name: "Past days targetPages ignored",
			progress: &models.ReadingProgress{
				TotalPages: 50,
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: yesterday, PagesRead: 0},
					{ID: 2, Date: today, PagesRead: 0},
					{ID: 3, Date: tomorrow, PagesRead: 0},
				},
			},
			logID:    2,
			expected: []int{10, 13, 13},
		},
		{
			name: "Future targetPages get evenly distributed",
			progress: &models.ReadingProgress{
				TotalPages: 30,
				EndDate:    today.AddDate(0, 0, 2),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: today, PagesRead: 0, TargetPages: 10},
					{ID: 2, Date: tomorrow, PagesRead: 0},
					{ID: 3, Date: today.AddDate(0, 0, 2), PagesRead: 0},
				},
			},
			logID:    0,
			expected: []int{10, 10, 10},
		},
		{
			name: "PagesRead 0 - simulating pagesRead equal target",
			progress: &models.ReadingProgress{
				TotalPages: 50,
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: yesterday, PagesRead: 10},
					{ID: 2, Date: today, PagesRead: 0, TargetPages: 10},
					{ID: 3, Date: tomorrow, PagesRead: 0, TargetPages: 10},
				},
			},
			logID:    2,
			expected: []int{10, 10, 10},
		},
		{
			name: "PagesRead equal TargetPages - same TP for next day",
			progress: &models.ReadingProgress{
				TotalPages: 50,
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: yesterday, PagesRead: 10},
					{ID: 2, Date: today, PagesRead: 10},
					{ID: 3, Date: tomorrow, PagesRead: 0},
				},
			},
			logID:    2,
			expected: []int{10, 10, 10},
		},
		{
			name: "PagesRead exceeds TargetPages -- smaller TP for next day",
			progress: &models.ReadingProgress{
				TotalPages: 50,
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: yesterday, PagesRead: 10},
					{ID: 2, Date: today, PagesRead: 25},
					{ID: 3, Date: tomorrow, PagesRead: 0},
				},
			},
			logID:    2,
			expected: []int{10, 10, 5},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := services.NewProgressService(nil)
			result, _ := service.UpdateTargetPages(tc.progress, tc.logID)

			var got []int
			for _, log := range result.DailyProgress {
				got = append(got, log.TargetPages)
			}

			assert.Equal(t, tc.expected, got, "Test failed for case: %s", tc.name)
		})
	}
}

func TestProgressService_UpdateTargetPagesForUserInput(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	progressService := services.NewProgressService(mockRepo)

	progressID := uint(123)
	progressIDString := "123"
	logID := uint(2)

	today := utils.TodaysDate()
	progress := &models.ReadingProgress{
		ID:         progressID,
		TotalPages: 50,
		EndDate:    today.AddDate(0, 0, 3),
		StartDate:  today.AddDate(0, 0, -1),
		Completed:  false,
		DailyProgress: []models.DailyProgressLog{
			{ID: 1, Date: today.AddDate(0, 0, -1), PagesRead: 10, TargetPages: 10},
			{ID: 2, Date: today, PagesRead: 0, TargetPages: 10},
			{ID: 3, Date: today.AddDate(0, 0, 1), PagesRead: 0, TargetPages: 10},
			{ID: 4, Date: today.AddDate(0, 0, 2), PagesRead: 0, TargetPages: 10},
			{ID: 5, Date: today.AddDate(0, 0, 3), PagesRead: 0, TargetPages: 10},
		},
	}

	t.Run("Normal case - TargetPages updated", func(t *testing.T) {
		mockRepo.On("Get", progressID).Return(progress, nil)
		mockRepo.On("GetById", progressIDString).Return(progress, nil)
		mockRepo.On("Update", progress).Return(nil)

		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.NoError(t, err)
		assert.Equal(t, progress, updatedProgress)

		assert.Equal(t, 10, progress.DailyProgress[0].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[1].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[2].TargetPages)
	})

	t.Run("Completed progress - No update", func(t *testing.T) {
		progress.Completed = true

		mockRepo.On("Get", progressID).Return(progress, nil)
		mockRepo.On("GetById", progressIDString).Return(progress, nil)

		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.NoError(t, err)
		assert.Equal(t, progress, updatedProgress)

		assert.Equal(t, 10, progress.DailyProgress[0].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[1].TargetPages)
		assert.Equal(t, 10, progress.DailyProgress[2].TargetPages)
	})
}

func TestProgressService_UpdateTargetPages_Invalid(t *testing.T) {
	mockRepo := new(MockProgressRepository)
	progressService := services.NewProgressService(mockRepo)
	progressIDString := "123"
	logID := uint(2)

	t.Run("Repo error - Get returns nil", func(t *testing.T) {
		mockRepo.On("GetById", progressIDString).Return(&models.ReadingProgress{}, errors.New("not found"))
		updatedProgress, err := progressService.UpdateTargetPagesForUserInput(progressIDString, logID)

		assert.Error(t, err)
		assert.Nil(t, updatedProgress)
		assert.Equal(t, "not found", err.Error())
	})
}
