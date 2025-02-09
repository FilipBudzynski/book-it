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
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestExchangeHandler(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	setupSession(t)

	e := setupEchoServerForExchange(t, db)

	user := &models.User{
		GoogleId: "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}
	require.NoError(t, db.Create(user).Error, "Failed to create user")

	userSession := utils.UserSession{
		UserID:    user.GoogleId,
		UserEmail: user.Email,
	}

	setSession := func(req *http.Request, userSession utils.UserSession) {
		rec := httptest.NewRecorder()
		err := utils.SetUserSession(rec, req, userSession)
		require.NoError(t, err, "Failed to set session")
	}

	t.Run("POST /exchange - Create Exchange", func(t *testing.T) {
		book := &models.Book{ID: "1", Title: "Test Book"}
		offeredBook1 := &models.Book{ID: "2", Title: "Offered Book 1"}
		offeredBook2 := &models.Book{ID: "3", Title: "Offered Book 2"}
		db.Create(book)
		db.Create(offeredBook1)
		db.Create(offeredBook2)

		formData := "latitude=50.06&longitude=19.94&desired-book-id=1&offered-book-0=2&offered-book-1=3"
		req := httptest.NewRequest(http.MethodPost, "/exchange", strings.NewReader(formData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange - Landing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange/list - Get All Exchanges", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange/list", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange/list/:status - Get Exchanges by Status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange/list/pending", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange/details/:id - Get Exchange Details", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange/details/1", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange/:id/matches - Get Exchange Matches", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange/1/matches", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /exchange/:id/matches/filter - Filtered Matches", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/exchange/1/matches/filter?distance=10", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("POST /exchange/accept/:id/:requestID - Accept Match", func(t *testing.T) {
		exchange := &models.ExchangeRequest{
			ID:            2,
			UserGoogleId:  "test-user-id",
			DesiredBookID: "1",
			OfferedBooks: []models.OfferedBook{
				{ID: 1, BookID: "2"},
				{ID: 2, BookID: "3"},
			},
			Status: "pending",
		}
		db.Create(exchange)

		matchRequest := &models.ExchangeMatch{
			ID:                       1,
			ExchangeRequestID:        1,
			MatchedExchangeRequestID: 2,
			Status:                   "pending",
			Request1Decision:         models.MatchDecisionPending,
			Request2Decision:         models.MatchDecisionPending,
		}
		db.Create(matchRequest)

		req := httptest.NewRequest(http.MethodPost, "/exchange/accept/1/2", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("POST /exchange/decline/:id/:requestID - Decline Match", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/exchange/decline/1/2", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE /exchange/:id - Delete Exchange", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/exchange/1", nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("GET /exchange/localization - Localization Autocomplete", func(t *testing.T) {
		query := "Plac Politechniki"
		req := httptest.NewRequest(http.MethodGet, "/exchange/localization?geoloc-query="+url.QueryEscape(query), nil)
		setSession(req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupEchoServerForExchange(t *testing.T, db *gorm.DB) *echo.Echo {
	t.Helper()
	e := echo.New()

	exchangeRepo := repositories.NewExchangeRequestRepository(db)
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)

	exchangeService := services.NewExchangeService(exchangeRepo)
	bookService := services.NewBookService(bookRepo).WithProvider(providers.NewGoogleProvider())
	userService := services.NewUserService(userRepo)

	exchangeHandler := handlers.NewExchangeHandler(exchangeService, bookService, userService)
	exchangeHandler.RegisterRoutes(e)

	return e
}
