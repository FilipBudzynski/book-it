package unit

import (
	"fmt"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestUserBookRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repositories.NewUserBookRepository(db)
	_, _, userBook := seedTestData(t, db, repo)

	t.Run("Create UserBook", func(t *testing.T) {
		err := db.Create(&models.Book{ID: "book456", Title: "Test Book 456"}).Error
		assert.NoError(t, err)

		newUserBook := &models.UserBook{UserGoogleId: "user123", BookID: "book456"}
		err = repo.Create(newUserBook)
		assert.NoError(t, err)

		got, err := repo.Get(fmt.Sprintf("%d", newUserBook.ID))
		assert.NoError(t, err)
		assert.Equal(t, "user123", got.UserGoogleId)
	})

	t.Run("Get UserBook", func(t *testing.T) {
		got, err := repo.Get(fmt.Sprintf("%d", userBook.ID))
		assert.NoError(t, err)
		assert.Equal(t, userBook.UserGoogleId, got.UserGoogleId)
	})

	t.Run("GetAllUserBooks", func(t *testing.T) {
		got, err := repo.GetAllUserBooks("user123")
		assert.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, userBook.BookID, got[0].BookID)
	})

	t.Run("Delete UserBook", func(t *testing.T) {
		err := repo.Delete(fmt.Sprintf("%d", userBook.ID))
		assert.NoError(t, err)

		got, err := repo.Get(fmt.Sprintf("%d", userBook.ID))
		assert.Error(t, err, "Record should not exist")
		assert.Equal(t, got, &models.UserBook{})
	})

	t.Run("Search UserBooks", func(t *testing.T) {
		results, err := repo.Search("user123", "Test")
		assert.NoError(t, err)
		assert.Len(t, results, 1)
	})
}

func seedTestData(t *testing.T, db *gorm.DB, repo services.UserBookRepository) (*models.User, *models.Book, *models.UserBook) {
	t.Helper()
	user := &models.User{
		GoogleId: "user123",
		Username: "testuser",
		Email:    "test@example.com",
	}
	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create user")

	book := &models.Book{
		ID:    "123",
		Title: "Test Book",
	}
	err = db.Create(book).Error
	require.NoError(t, err, "Failed to create test book")

	userBook := &models.UserBook{
		UserGoogleId: "user123",
		BookID:       "123",
	}
	err = repo.Create(userBook)
	require.NoError(t, err, "Failed to create user book")

	return user, book, userBook
}
