package integration_tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserBookHandler(t *testing.T) {
	db, clean := setupTestDB(t)
	defer clean()
	setupSession(t)

	e := setupUserBookHandler(db)
	user, book, userBook := seedUserBookHandlerData(t, db)

	userSession := utils.UserSession{
		UserID:    user.GoogleId,
		UserEmail: user.Email,
	}

	t.Run("POST /user-books/:book_id - Create User Book", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/user-books/"+book.ID, nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE /user-books/:book_id - Delete User Book", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user-books/%d", userBook.ID), nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /user-books - List User Books", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user-books", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /user-books/create_modal/:user_book_id - Get Create Progress Modal", func(t *testing.T) {
		userBook := models.UserBook{}
		db.Where("user_google_id = ?", user.GoogleId).First(&userBook)

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user-books/create_modal/%d", userBook.ID), nil)
		setSession(t, req, utils.UserSession{
			UserID:    user.GoogleId,
			UserEmail: user.Email,
		})

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d, Response Body: %s", rec.Code, rec.Body.String())

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupUserBookHandler(db *gorm.DB) *echo.Echo {
	e := echo.New()

	userBookRepo := repositories.NewUserBookRepository(db)
	exchangeRequestRepo := repositories.NewExchangeRequestRepository(db)

	userBookService := services.NewUserBookService(userBookRepo, exchangeRequestRepo)

	userBookHandler := handlers.NewUserBookHandler(userBookService)
	userBookHandler.RegisterRoutes(e)
	return e
}

func seedUserBookHandlerData(t *testing.T, db *gorm.DB) (*models.User, *models.Book, *models.UserBook) {
	t.Helper()
	user := &models.User{
		GoogleId: "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	book := &models.Book{
		ID:    "1",
		Title: "Test Book",
	}
	if err := db.Create(book).Error; err != nil {
		t.Fatalf("Failed to create book: %v", err)
	}

	userBook := &models.UserBook{
		BookID:       book.ID,
		UserGoogleId: user.GoogleId,
	}
	if err := db.Create(userBook).Error; err != nil {
		t.Fatalf("Failed to create user book: %v", err)
	}
	return user, book, userBook
}
