package unit

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserBookRepository struct {
	mock.Mock
}

func (m *MockUserBookRepository) Create(userBook *models.UserBook) error {
	args := m.Called(userBook)
	return args.Error(0)
}

func (m *MockUserBookRepository) Update(userBook *models.UserBook) error {
	args := m.Called(userBook)
	return args.Error(0)
}

func (m *MockUserBookRepository) Get(id string) (*models.UserBook, error) {
	args := m.Called(id)
	return args.Get(0).(*models.UserBook), args.Error(1)
}

func (m *MockUserBookRepository) GetAllUserBooks(userId string) ([]*models.UserBook, error) {
	args := m.Called(userId)
	return args.Get(0).([]*models.UserBook), args.Error(1)
}

func (m *MockUserBookRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserBookRepository) DeleteWhereBookId(bookId string) error {
	args := m.Called(bookId)
	return args.Error(0)
}

func (m *MockUserBookRepository) Search(userId, query string) ([]*models.UserBook, error) {
	args := m.Called(userId, query)
	return args.Get(0).([]*models.UserBook), args.Error(1)
}

