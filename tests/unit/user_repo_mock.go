package unit

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetById(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	args := m.Called(googleID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(userID string) error {
	args := m.Called(userID)
	return args.Error(0)
}

func (m *MockUserRepository) FirstGenre(genreID string) (*models.Genre, error) {
	args := m.Called(genreID)
	return args.Get(0).(*models.Genre), args.Error(1)
}

func (m *MockUserRepository) AddGenre(user *models.User, genre *models.Genre) error {
	args := m.Called(user, genre)
	return args.Error(0)
}

func (m *MockUserRepository) RemoveGenre(user *models.User, genre *models.Genre) error {
	args := m.Called(user, genre)
	return args.Error(0)
}

func (m *MockUserRepository) GetAllGenres() ([]*models.Genre, error) {
	args := m.Called()
	return args.Get(0).([]*models.Genre), args.Error(1)
}

func (m *MockUserRepository) ClearExpectedCalls() {
	m.ExpectedCalls = nil
}
