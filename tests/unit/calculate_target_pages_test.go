package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestCalculateTargetPages(t *testing.T) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := services.CalculateTargetPages(tt.pagesLeft, tt.daysLeft)
			assert.Equal(t, tt.expected, result, "Expected %d but got %d", tt.expected, result)
		})
	}
}
