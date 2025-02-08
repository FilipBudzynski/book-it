package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func TestProgressHandler(t *testing.T) {
	db, clean := setupTestDB(t)
	defer clean()
	setupSession(t)

	e := setupProgressHandler(t, db)

	user, _, _ := seedProgresHandlerData(t, db)

	userSession := utils.UserSession{
		UserID:    user.GoogleId,
		UserEmail: user.Email,
	}

	t.Run("POST /progress - Create Progress", func(t *testing.T) {
		form := url.Values{}
		form.Add("user-book-id", "1")
		form.Add("total-pages", "350")
		form.Add("book-title", "Test Book")
		form.Add("start-date", "2025-01-01")
		form.Add("end-date", "2025-02-01")

		req := httptest.NewRequest(http.MethodPost, "/progress", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /progress/:id - Get Progress by User Book ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/progress/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /progress/details/:id - Get Progress Details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/progress/details/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("PUT /progress/log/:id - Update Log", func(t *testing.T) {
		form := url.Values{}
		form.Add("comment", "Read 20 pages")
		form.Add("pages-read", "20")

		req := httptest.NewRequest(http.MethodPut, "/progress/log/1", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /progress/log/details/modal/:id - Get Log Modal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/progress/log/details/modal/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE /progress/:id - Delete Progress", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/progress/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupProgressHandler(t *testing.T, db *gorm.DB) *echo.Echo {
	t.Helper()
	t.Setenv("test-var", "true")
	e := echo.New()

	progressRepo := repositories.NewProgressRepository(db)
	userBookRepo := repositories.NewUserBookRepository(db)
	exchangeRequestRepo := repositories.NewExchangeRequestRepository(db)

	progressService := services.NewProgressService(progressRepo)
	userBookService := services.NewUserBookService(userBookRepo, exchangeRequestRepo)

	progressHandler := handlers.NewProgressHandler(progressService, userBookService)
	progressHandler.RegisterRoutes(e)
	return e
}

func seedProgresHandlerData(t *testing.T, db *gorm.DB) (*models.User, *models.Book, *models.UserBook) {
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
