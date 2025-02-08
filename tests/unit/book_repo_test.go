package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestBookRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := repositories.NewBookRepository(db)

	t.Run("Create Book", func(t *testing.T) {
		book := &models.Book{
			ID:    "1",
			Title: "Test Book",
			Genres: []models.Genre{
				{Name: "Fiction"},
			},
		}
		err := repo.Create(book)
		assert.NoError(t, err)
	})

	t.Run("Get Book By ID", func(t *testing.T) {
		book, err := repo.Get("1")
		assert.NoError(t, err)
		assert.NotNil(t, book)
		assert.Equal(t, "Test Book", book.Title)
	})

	t.Run("Delete Book", func(t *testing.T) {
		err := repo.Delete("1")
		assert.NoError(t, err)

		book, err := repo.Get("1")
		assert.Error(t, err)
		assert.Nil(t, book)
	})

	t.Run("Get Books By Genre", func(t *testing.T) {
		book1 := &models.Book{ID: "2", Title: "Sci-Fi Book", Genres: []models.Genre{{Name: "Sci-Fi"}}}
		book2 := &models.Book{ID: "3", Title: "Another Sci-Fi Book", Genres: []models.Genre{{Name: "Sci-Fi"}}}
		repo.Create(book1)
		repo.Create(book2)

		books, err := repo.GetByGenre("Sci-Fi")
		assert.NoError(t, err)
		assert.Len(t, books, 2)
	})
}
