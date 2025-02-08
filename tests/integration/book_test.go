package integration_tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/providers"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookHandler(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	setupSession(t)

	e := setupEchoServer(t, db)

	user := &models.User{
		GoogleId: "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	userSession := utils.UserSession{
		UserID:    user.GoogleId,
		UserEmail: user.Email,
	}

	setSession := func(req *http.Request, userSession utils.UserSession) {
		rec := httptest.NewRecorder()
		err := utils.SetUserSession(rec, req, userSession)
		if err != nil {
			t.Fatalf("Failed to set session: %v", err)
		}
	}

	t.Run("GET /books", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("POST /books", func(t *testing.T) {
		formData := "query=somebook&type=title"
		req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(formData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /books/reduced/search", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/reduced/search", nil)
		req.Form = url.Values{
			"book-title": {"somebook"},
		}

		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /books/partial", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/partial", nil)
		q := req.URL.Query()
		q.Add("query", "somebook")
		q.Add("page", "1")
		req.URL.RawQuery = q.Encode()

		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("GET /books/recommendations", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books/recommendations", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupEchoServer(t *testing.T, db *gorm.DB) *echo.Echo {
	t.Helper()
	t.Setenv("test-var", "true")
	e := echo.New()

	userBookRepo := repositories.NewUserBookRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	userRepo := repositories.NewUserRepository(db)
	exchangeRequestRepo := repositories.NewExchangeRequestRepository(db)

	bookService := services.NewBookService(bookRepo).WithProvider(providers.NewGoogleProvider())
	userBookService := services.NewUserBookService(userBookRepo, exchangeRequestRepo)
	userService := services.NewUserService(userRepo)

	bookHandler := handlers.NewBookHandler(bookService, userBookService, userService)
	bookHandler.RegisterRoutes(e)

	return e
}
