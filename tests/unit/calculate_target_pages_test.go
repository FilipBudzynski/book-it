package unit

import (
	"testing"
	"time"

	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTargetPages(t *testing.T) {
	tests := []struct {
		name      string
		pagesLeft int
		daysLeft  int
		logDate   time.Time
		expected  int
	}{
		{"Zero Days Left", 100, 0, time.Now(), 100},

		{"Even Distribution", 100, 10, time.Now(), 10},

		{"Rounding Up", 101, 10, time.Now(), 11},

		{"No Pages Left", 0, 10, time.Now(), 0},
		{"Negative Pages Left", -5, 10, time.Now(), -1},

		{"Negative Days Left", 100, -10, time.Now(), -1},

		{"One Day Left", 50, 1, time.Now(), 50},
		{"One Page Left, One Day Left", 1, 1, time.Now(), 1},
		{"One Page Left, Multiple Days", 1, 10, time.Now(), 1},

		{"Large Pages and Days", 1000000, 1000, time.Now(), 1000},
		{"Large Pages, Few Days", 1000000, 10, time.Now(), 100000},
		{"Few Pages, Large Days", 10, 1000, time.Now(), 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.CalculateTargetPages(tt.pagesLeft, tt.daysLeft, tt.logDate)
			assert.Equal(t, tt.expected, result, "Expected %d but got %d", tt.expected, result)
		})
	}
}
