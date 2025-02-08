package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestDeleteUserBook(t *testing.T) {
	userBookID := "user-book-1"
	bookID := "book-123"
	userID := "user-456"

	userBook := &models.UserBook{
		UserGoogleId: userID,
		BookID:       bookID,
	}

	activeRequests := []*models.ExchangeRequest{
		{
			ID: 123,
		},
	}

	mockUserBookRepo := new(MockUserBookRepository)
	mockUserBookRepo.On("Get", userBookID).Return(userBook, nil)

	mockExchangeRepo := new(MockExchangeRequestRepository)
	mockExchangeRepo.On("GetActiveExchangeRequestsByBookID", bookID, userID).Return(activeRequests, nil)

	service := services.NewUserBookService(mockUserBookRepo, mockExchangeRepo)

	err := service.Delete(userBookID)

	assert.NotNil(t, err)
	assert.Equal(t, models.ErrUserBookInActiveExchangeRequest, err)

	mockUserBookRepo.AssertExpectations(t)
	mockExchangeRepo.AssertExpectations(t)
}

func TestDeleteUserBook_NoActiveRequests(t *testing.T) {
	userBookID := "user-book-1"
	bookID := "book-123"
	userID := "user-456"

	userBook := &models.UserBook{
		UserGoogleId: userID,
		BookID:       bookID,
	}

	mockExchangeRepo := new(MockExchangeRequestRepository)
	mockExchangeRepo.On("GetActiveExchangeRequestsByBookID", bookID, userID).Return([]*models.ExchangeRequest{}, nil)

	mockUserBookRepo := new(MockUserBookRepository)
	mockUserBookRepo.On("Get", userBookID).Return(userBook, nil)
	mockUserBookRepo.On("Delete", userBookID).Return(nil)

	service := services.NewUserBookService(mockUserBookRepo, mockExchangeRepo)

	err := service.Delete(userBookID)

	assert.Nil(t, err)

	mockUserBookRepo.AssertExpectations(t)
	mockExchangeRepo.AssertExpectations(t)
}
