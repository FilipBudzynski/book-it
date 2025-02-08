package unit

import (
	"errors"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := services.NewUserService(mockRepo)

	t.Run("Create - success", func(t *testing.T) {
		user := &models.User{Email: "test@example.com", Username: "testuser", GoogleId: "123"}
		mockRepo.On("Create", user).Return(nil)

		err := service.Create(user)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Create", user)
	})

	t.Run("Create - validation error", func(t *testing.T) {
		user := &models.User{}
		mockRepo.On("Create", user).Return(nil)

		err := service.Create(user)

		assert.Error(t, err)
		mockRepo.AssertNotCalled(t, "Create", user)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("AddGenre - success", func(t *testing.T) {
		user := &models.User{GoogleId: "user123", Genres: []models.Genre{}}
		genre := &models.Genre{Name: "Fantasy"}

		mockRepo.On("GetByGoogleID", "user123").Return(user, nil)
		mockRepo.On("FirstGenre", "genre123").Return(genre, nil)
		mockRepo.On("AddGenre", user, genre).Return(nil)

		result, err := service.AddGenre("user123", "genre123")

		assert.NoError(t, err)
		assert.Equal(t, genre, result)
		mockRepo.AssertCalled(t, "AddGenre", user, genre)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("AddGenre - genre limit exceeded", func(t *testing.T) {
		user := &models.User{
			GoogleId: "user123",
			Genres: []models.Genre{
				{ID: 1, Name: "Fantasy"},
				{ID: 2, Name: "Sci-Fi"},
				{ID: 3, Name: "Horror"},
				{ID: 4, Name: "Mystery"},
				{ID: 5, Name: "Romance"},
			},
		}

		mockRepo.On("GetByGoogleID", "user123").Return(user, nil)

		result, err := service.AddGenre("user123", "1")

		assert.ErrorIs(t, err, models.ErrUserGenresLimitExceeded)
		assert.Nil(t, result)

		mockRepo.AssertNotCalled(t, "FirstGenre")
		mockRepo.AssertNotCalled(t, "AddGenre")
		mockRepo.ClearExpectedCalls()
	})

	t.Run("RemoveGenre - success", func(t *testing.T) {
		user := &models.User{GoogleId: "user123"}
		genre := &models.Genre{Name: "Fantasy"}

		mockRepo.On("GetByGoogleID", "user123").Return(user, nil)
		mockRepo.On("FirstGenre", "genre123").Return(genre, nil)
		mockRepo.On("RemoveGenre", user, genre).Return(nil)

		result, err := service.RemoveGenre("user123", "genre123")

		assert.NoError(t, err)
		assert.Equal(t, genre, result)
		mockRepo.AssertCalled(t, "RemoveGenre", user, genre)
		mockRepo.ClearExpectedCalls()
	})

	t.Run("RemoveGenre - user not found", func(t *testing.T) {
		mockRepo.On("GetByGoogleID", "user123").Return(&models.User{}, errors.New("not found"))

		result, err := service.RemoveGenre("user123", "genre123")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "RemoveGenre")
		mockRepo.ClearExpectedCalls()
	})

	t.Run("RemoveGenre - genre not found", func(t *testing.T) {
		user := &models.User{GoogleId: "user123"}

		mockRepo.On("GetByGoogleID", "user123").Return(user, nil)
		mockRepo.On("FirstGenre", "genre000").Return(&models.Genre{}, errors.New("not found"))

		result, err := service.RemoveGenre("user123", "genre000")

		assert.Error(t, err)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "RemoveGenre")
		mockRepo.ClearExpectedCalls()
	})

	t.Run("GetByGoogleID - success", func(t *testing.T) {
		user := &models.User{GoogleId: "user123"}
		mockRepo.On("GetByGoogleID", "user123").Return(user, nil)

		result, err := service.GetByGoogleID("user123")

		assert.NoError(t, err)
		assert.Equal(t, user, result)
		mockRepo.AssertCalled(t, "GetByGoogleID", "user123")
	})

	t.Run("GetAll - success", func(t *testing.T) {
		users := []models.User{{GoogleId: "user1"}, {GoogleId: "user2"}}
		mockRepo.On("GetAll").Return(users, nil)

		result, err := service.GetAll()

		assert.NoError(t, err)
		assert.Equal(t, users, result)
		mockRepo.AssertCalled(t, "GetAll")
	})

	t.Run("Delete - success", func(t *testing.T) {
		mockRepo.On("Delete", "user123").Return(nil)

		err := service.Delete("user123")

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "Delete", "user123")
	})
}
