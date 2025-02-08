package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)
	user, genre := seedTestUserData(t, db)

	t.Run("Create User", func(t *testing.T) {
		newUser := &models.User{
			GoogleId: "google456",
			Username: "testuser456",
			Email:    "testuser456@example.com",
		}
		err := repo.Create(newUser)
		assert.NoError(t, err)

		got, err := repo.GetByGoogleID("google456")
		assert.NoError(t, err)
		assert.Equal(t, "testuser456", got.Username)
	})

	t.Run("Get User by Google ID", func(t *testing.T) {
		got, err := repo.GetByGoogleID(user.GoogleId)
		assert.NoError(t, err)
		assert.Equal(t, user.Username, got.Username)
	})

	t.Run("Get User by Email", func(t *testing.T) {
		got, err := repo.GetByEmail(user.Email)
		assert.NoError(t, err)
		assert.Equal(t, user.Email, got.Email)
	})

	t.Run("Get All Users", func(t *testing.T) {
		users, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, user.GoogleId, users[0].GoogleId)
		assert.Equal(t, "google456", users[1].GoogleId)
	})

	t.Run("Update User", func(t *testing.T) {
		user.Username = "updatedusername"
		err := repo.Update(user)
		assert.NoError(t, err)

		updatedUser, err := repo.GetByGoogleID(user.GoogleId)
		assert.NoError(t, err)
		assert.Equal(t, "updatedusername", updatedUser.Username)
	})

	t.Run("Add Genre to User", func(t *testing.T) {
		err := repo.AddGenre(user, genre)
		assert.NoError(t, err)

		updatedUser, err := repo.GetByGoogleID(user.GoogleId)
		assert.NoError(t, err)
		assert.Equal(t, updatedUser.Genres[0].Name, genre.Name)
	})

	t.Run("Remove Genre from User", func(t *testing.T) {
		err := repo.RemoveGenre(user, genre)
		assert.NoError(t, err)

		got, err := repo.GetByGoogleID(user.GoogleId)
		assert.NoError(t, err)
		assert.Len(t, got.Genres, 0)
	})

	t.Run("Delete User", func(t *testing.T) {
		err := repo.Delete(user.GoogleId)
		assert.NoError(t, err)

		got, err := repo.GetByGoogleID(user.GoogleId)
		assert.Error(t, err)
		assert.Equal(t, got, &models.User{})
	})
}

func seedTestUserData(t *testing.T, db *gorm.DB) (*models.User, *models.Genre) {
	t.Helper()

	user := &models.User{
		GoogleId: "google123",
		Username: "testuser",
		Email:    "testuser@example.com",
	}
	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create user")

	genre := &models.Genre{Name: "Science Fiction"}
	err = db.Create(genre).Error
	require.NoError(t, err, "Failed to create genre")

	return user, genre
}
