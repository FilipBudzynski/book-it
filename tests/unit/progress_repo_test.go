package unit

import (
	"fmt"
	"testing"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestProgressRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewProgressRepository(db)
	_, _, userBook := seedProgressTestData(t, db)
	progressID := "123"

	t.Run("Create Reading Progress", func(t *testing.T) {
		newProgress := models.ReadingProgress{
			ID:         123,
			UserBookID: userBook.ID,
			TotalPages: 200,
			DailyProgress: []models.DailyProgressLog{
				{
					Date:      time.Now(),
					PagesRead: 10,
				},
			},
		}
		err := repo.Create(newProgress)
		assert.NoError(t, err)

		got, err := repo.GetById(fmt.Sprintf("%d", newProgress.ID))
		assert.NoError(t, err)
		assert.Equal(t, userBook.ID, got.UserBookID)
	})

	t.Run("GetById", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		got, err := repo.GetById(fmt.Sprintf("%d", progress.ID))
		assert.NoError(t, err)
		assert.Equal(t, progress, got)
	})

	t.Run("GetByUserBookId", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		got, err := repo.GetByUserBookId(fmt.Sprintf("%d", userBook.ID))
		assert.NoError(t, err)
		assert.Equal(t, progress, got)
	})

	t.Run("GetLogById", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		got, err := repo.GetLogById(fmt.Sprintf("%d", progress.DailyProgress[0].ID))
		assert.NoError(t, err)
		assert.Equal(t, progress.DailyProgress[0].PagesRead, got.PagesRead)
	})

	t.Run("Update Reading Progress", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		progress.TotalPages = 300
		err = repo.Update(progress)
		assert.NoError(t, err)

		updated, err := repo.GetById(fmt.Sprintf("%d", progress.ID))
		assert.NoError(t, err)
		assert.Equal(t, 300, updated.TotalPages)
	})

	t.Run("Update Daily Progress Log", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		log := progress.DailyProgress[0]
		log.PagesRead = 15
		err = repo.UpdateLog(&log)
		assert.NoError(t, err)

		updated, _ := repo.GetLogById(fmt.Sprintf("%d", log.ID))
		assert.Equal(t, 15, updated.PagesRead)
	})

	t.Run("Delete Reading Progress", func(t *testing.T) {
		progress, err := repo.GetById(progressID)
		assert.NoError(t, err)

		err = repo.Delete(fmt.Sprintf("%d", progress.ID))
		assert.NoError(t, err)

		got, err := repo.GetById(fmt.Sprintf("%d", progress.ID))
		assert.Error(t, err, "Record should not exist")
		assert.Nil(t, got)
	})
}

func seedProgressTestData(t *testing.T, db *gorm.DB) (*models.User, *models.Book, *models.UserBook) {
	t.Helper()

	user := &models.User{
		GoogleId: "user123",
		Username: "testuser",
		Email:    "test@example.com",
	}
	require.NoError(t, db.Create(user).Error, "Failed to create user")

	book := &models.Book{
		ID:    "book123",
		Title: "Test Book",
	}
	require.NoError(t, db.Create(book).Error, "Failed to create book")

	userBook := &models.UserBook{
		UserGoogleId: user.GoogleId,
		BookID:       book.ID,
	}
	require.NoError(t, db.Create(userBook).Error, "Failed to create user book")
	return user, book, userBook
}
