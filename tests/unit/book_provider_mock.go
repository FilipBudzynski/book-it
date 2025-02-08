package unit

import (
	"github.com/FilipBudzynski/book_it/internal/handlers"
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBookProvider struct {
	mock.Mock
}

func (m *MockBookProvider) GetBook(id string) (*models.Book, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *MockBookProvider) GetBooksByQuery(query string, queryType handlers.QueryType, limit, page int) ([]*models.Book, error) {
	args := m.Called(query, queryType, limit, page)
	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookProvider) GetBooksByGenre(genre string) ([]*models.Book, error) {
	args := m.Called(genre)
	return args.Get(0).([]*models.Book), args.Error(1)
}

func (m *MockBookProvider) QueryTypeToString(queryType handlers.QueryType) string {
	args := m.Called(queryType)
	return args.String(0)
}

func (m *MockBookProvider) Convert(response any) *models.Book {
	args := m.Called(response)
	return args.Get(0).(*models.Book)
}
