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

func TestUserHandler(t *testing.T) {
	db, clean := setupTestDB(t)
	defer clean()
	setupSession(t)

	e := setupUserHandler(t, db)

	user, _ := seedUserHandlerData(t, db)

	userSession := utils.UserSession{
		UserID:    user.GoogleId,
		UserEmail: user.Email,
	}

	t.Run("GET /users/profile - Fetch Profile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/profile", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d, Body: %s", rec.Code, rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("GET /users/profile/location/modal - Get Location Modal", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/profile/location/modal", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d, Body: %s", rec.Code, rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("POST /users/profile/genres/:genre_id - Add Genre", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/profile/genres/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d, Body: %s", rec.Code, rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE /users/profile/genres/:genre_id - Remove Genre", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/profile/genres/1", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d, Body: %s", rec.Code, rec.Body.String())
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("POST /users/profile/location - Change Location", func(t *testing.T) {
		form := url.Values{}
		form.Add("latitude", "40.7128")
		form.Add("longitude", "-74.0060")
		form.Add("formatted", "New York, USA")

		req := httptest.NewRequest(http.MethodPost, "/users/profile/location", strings.NewReader(form.Encode()))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d", rec.Code)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("DELETE /users - Delete User Account", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users", nil)
		setSession(t, req, userSession)

		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		t.Logf("Response Code: %d", rec.Code)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func setupUserHandler(t *testing.T, db *gorm.DB) *echo.Echo {
	t.Helper()
	t.Setenv("test-var", "true")
	e := echo.New()

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService)
	userHandler.RegisterRoutes(e)
	return e
}

func seedUserHandlerData(t *testing.T, db *gorm.DB) (*models.User, *models.Genre) {
	t.Helper()
	user := &models.User{
		GoogleId: "test-user-id",
		Username: "testuser",
		Email:    "test@example.com",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	genre := &models.Genre{
		Name: "Science Fiction",
	}
	if err := db.Create(genre).Error; err != nil {
		t.Fatalf("Failed to create genre: %v", err)
	}

	return user, genre
}
