package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTargetPages(t *testing.T) {
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
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: tomorrow, PagesRead: 0},
					{ID: 2, Date: today.AddDate(0, 0, 2), PagesRead: 0},
				},
			},
			logID:    0,
			expected: []int{10, 10}, 
		},
		{
			name: "PagesRead 0 - simulating pg equal tp",
			progress: &models.ReadingProgress{
				TotalPages: 50,
				EndDate:    today.AddDate(0, 0, 3),
				DailyProgress: []models.DailyProgressLog{
					{ID: 1, Date: yesterday, PagesRead: 10},
					{ID: 2, Date: today, PagesRead: 0}, 
					{ID: 3, Date: tomorrow, PagesRead: 0},
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
