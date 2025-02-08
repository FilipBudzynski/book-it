package unit

import (
	"fmt"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestExchangeRequestRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := repositories.NewExchangeRequestRepository(db)
	user, _, _ := seedExchangeRequestTestData(t, db)
	exchange := &models.ExchangeRequest{}

	t.Run("Create Exchange Request", func(t *testing.T) {
		newExchange := &models.ExchangeRequest{
			UserGoogleId:  user.GoogleId,
			DesiredBookID: "book123",
			Status:        models.ExchangeRequestStatusActive,
		}

		err := repo.Create(newExchange)
		assert.NoError(t, err)

		got, err := repo.GetByID(fmt.Sprintf("%d", newExchange.ID))
		assert.NoError(t, err)
		assert.Equal(t, newExchange.UserGoogleId, got.UserGoogleId)
		exchange = got
	})

	t.Run("GetByID", func(t *testing.T) {
		got, err := repo.GetByID(fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)
		assert.Equal(t, exchange, got)
	})

	t.Run("GetByUserId", func(t *testing.T) {
		exchange, err := repo.GetByID("1")
		assert.NoError(t, err)

		got, err := repo.Get(fmt.Sprintf("%d", exchange.ID), exchange.UserGoogleId)
		assert.NoError(t, err)
		assert.Equal(t, exchange, got)
	})

	t.Run("GetAll", func(t *testing.T) {
		exchanges, err := repo.GetAll(user.GoogleId)
		assert.NoError(t, err)
		assert.NotEmpty(t, exchanges)
	})

	t.Run("GetAllWithStatus", func(t *testing.T) {
		exchanges, err := repo.GetAllWithStatus(user.GoogleId, models.ExchangeRequestStatusActive)
		assert.NoError(t, err)
		assert.NotEmpty(t, exchanges)
	})

	t.Run("Delete Exchange Request", func(t *testing.T) {
		err := repo.Delete(fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)

		got, err := repo.GetByID(fmt.Sprintf("%d", exchange.ID))
		assert.Error(t, err, "Record should not exist")
		assert.Nil(t, got)
	})

	t.Run("Delete Matches for Request", func(t *testing.T) {
		err := repo.DeleteMatchesForRequest(fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)

		gotMatches, err := repo.GetAllMatches(fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)
		assert.Empty(t, gotMatches)
	})

	t.Run("Update Exchange Request", func(t *testing.T) {
		newExchange := &models.ExchangeRequest{
			UserGoogleId:  user.GoogleId,
			DesiredBookID: "book123",
			Status:        models.ExchangeRequestStatusActive,
		}
		err := repo.Create(newExchange)
		assert.NoError(t, err)
		exchange := newExchange

		exchange.Status = models.ExchangeRequestStatusCompleted
		err = repo.Update(exchange)
		assert.NoError(t, err)

		updated, err := repo.GetByID(fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)
		assert.Equal(t, models.ExchangeRequestStatusCompleted, updated.Status)
	})

	t.Run("Create Match", func(t *testing.T) {
		exchange, err := repo.GetByID("1")
		assert.NoError(t, err)

		match := &models.ExchangeMatch{
			ExchangeRequestID:        exchange.ID,
			MatchedExchangeRequestID: 2,
		}

		err = repo.CreateMatch(match)
		assert.NoError(t, err)

		gotMatch, err := repo.GetMatch(fmt.Sprintf("%d", match.ID), fmt.Sprintf("%d", exchange.ID))
		assert.NoError(t, err)
		assert.NotNil(t, gotMatch)
	})
}

func seedExchangeRequestTestData(t *testing.T, db *gorm.DB) (*models.User, *models.Book, *models.ExchangeRequest) {
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

	exchangeRequest := &models.ExchangeRequest{
		UserGoogleId:  user.GoogleId,
		DesiredBookID: book.ID,
		Status:        models.ExchangeRequestStatusActive,
	}
	require.NoError(t, db.Create(exchangeRequest).Error, "Failed to create exchange request")
	return user, book, exchangeRequest
}
