package unit

import (
	"errors"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestBookService(t *testing.T) {
	repo := new(MockBookRepository)
	provider := new(MockBookProvider)
	svc := services.NewBookService(repo).WithProvider(provider)

	t.Run("Create", func(t *testing.T) {
		book := &models.Book{ID: "1", Title: "Test Book"}
		repo.On("Create", book).Return(nil)
		err := svc.Create(book)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("GetByID - Exists in Repo", func(t *testing.T) {
		book := &models.Book{ID: "1", Title: "Test Book"}
		repo.On("Get", "1").Return(book, nil)
		got, err := svc.GetByID("1")
		assert.NoError(t, err)
		assert.Equal(t, book, got)
		repo.AssertExpectations(t)
	})

	t.Run("GetByID - Fetch from Provider", func(t *testing.T) {
		book := &models.Book{ID: "2", Title: "Provider Book"}
		repo.On("Get", "2").Return((*models.Book)(nil), errors.New("not found"))
		provider.On("GetBook", "2").Return(book, nil)
		repo.On("Create", book).Return(nil)
		got, err := svc.GetByID("2")
		assert.NoError(t, err)
		assert.Equal(t, book, got)
		repo.AssertExpectations(t)
		provider.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		repo.On("Delete", "1").Return(nil)
		err := svc.Delete("1")
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("GetByQuery", func(t *testing.T) {
		books := []*models.Book{{ID: "3", Title: "Query Book"}}
		provider.On("GetBooksByQuery", "test", handlers.QueryTypeTitle, 40, 1).Return(books, nil)
		repo.On("Get", "3").Return((*models.Book)(nil), errors.New("not found"))
		repo.On("Create", books[0]).Return(nil)
		got, err := svc.GetByQuery("test", handlers.QueryTypeTitle, 1)
		assert.NoError(t, err)
		assert.Equal(t, books, got)
		repo.AssertExpectations(t)
		provider.AssertExpectations(t)
	})

	t.Run("FetchReccomendations", func(t *testing.T) {
		genres := []models.Genre{{Name: "Sci-Fi"}}
		books := []*models.Book{{ID: "4", Title: "Sci-Fi Book"}}
		userBooks := []*models.UserBook{{Book: models.Book{ID: "5"}}}
		provider.On("GetBooksByGenre", "Sci-Fi").Return(books, nil)
		repo.On("Get", "4").Return((*models.Book)(nil), errors.New("not found"))
		repo.On("Create", books[0]).Return(nil)
		got, err := svc.FetchReccomendations(genres, userBooks)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		assert.Equal(t, books[0], got[0])
		repo.AssertExpectations(t)
		provider.AssertExpectations(t)
	})
}
