package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
)

func TestReadingProgressEqual(t *testing.T) {
	tests := []struct {
		name     string
		r1       models.ReadingProgress
		r2       models.ReadingProgress
		expected bool
	}{
		{
			name: "Equal ReadingProgress instances",
			r1: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			r2: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			expected: true,
		},
		{
			name: "Different DailyTargetPages",
			r1: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			r2: models.ReadingProgress{
				DailyTargetPages: 6,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			expected: false,
		},
		{
			name: "Different number of Log entries",
			r1: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
				},
			},
			r2: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			expected: false,
		},
		{
			name: "Different TargetPages in Logs",
			r1: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 3},
				},
			},
			r2: models.ReadingProgress{
				DailyTargetPages: 5,
				DailyProgress: []models.DailyProgressLog{
					{TargetPages: 2},
					{TargetPages: 4},
				},
			},
			expected: false,
		},
		{
			name:     "Empty ReadingProgress instances",
			r1:       models.ReadingProgress{},
			r2:       models.ReadingProgress{},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.r1.Equal(tt.r2)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
