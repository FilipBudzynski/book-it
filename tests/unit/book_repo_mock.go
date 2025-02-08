package unit

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) Get(id string) (*models.Book, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Book), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockBookRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockBookRepository) GetByGenre(genre string) ([]*models.Book, error) {
	args := m.Called(genre)
	if args.Get(0) != nil {
		return args.Get(0).([]*models.Book), args.Error(1)
	}
	return nil, args.Error(1)
}

