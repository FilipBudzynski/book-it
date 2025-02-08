package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name     string
		user     *models.User
		expected error
	}{
		{
			name: "Valid User",
			user: &models.User{
				Username: "john_doe",
				Email:    "john.doe@example.com",
				GoogleId: "google123",
			},
			expected: nil,
		},
		{
			name: "Empty Username",
			user: &models.User{
				Username: "",
				Email:    "john.doe@example.com",
				GoogleId: "google123",
			},
			expected: models.ErrUsernameRequired,
		},
		{
			name: "Invalid Email",
			user: &models.User{
				Username: "john_doe",
				Email:    "john.doe",
				GoogleId: "google123",
			},
			expected: models.ErrEmailRequired,
		},
		{
			name: "Empty GoogleId",
			user: &models.User{
				Username: "john_doe",
				Email:    "john.doe@example.com",
				GoogleId: "",
			},
			expected: models.ErrGoogleIdRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.expected != nil {
				assert.Equal(t, tt.expected, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_HasGenre(t *testing.T) {
	genre1 := models.Genre{Name: "Fantasy"}
	genre2 := models.Genre{Name: "Sci-Fi"}

	tests := []struct {
		name     string
		user     *models.User
		genre    string
		expected bool
	}{
		{
			name:     "Has genre",
			user:     &models.User{Genres: []models.Genre{genre1, genre2}},
			genre:    "Fantasy",
			expected: true,
		},
		{
			name:     "Does not have genre",
			user:     &models.User{Genres: []models.Genre{genre1, genre2}},
			genre:    "Horror",
			expected: false,
		},
		{
			name:     "Empty genre list",
			user:     &models.User{Genres: []models.Genre{}},
			genre:    "Fantasy",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.HasGenre(tt.genre)
			assert.Equal(t, tt.expected, result)
		})
	}
}
