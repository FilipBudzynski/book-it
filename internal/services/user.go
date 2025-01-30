package services

import (
	"errors"
	"fmt"

	"github.com/FilipBudzynski/book_it/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetById(id string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(userID string) error
	FirstGenre(genreID string) (*models.Genre, error)
	GetAllGenres() ([]*models.Genre, error)
	AddGenre(user *models.User, genre *models.Genre) error
	RemoveGenre(user *models.User, genre *models.Genre) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *userService {
	return &userService{
		repo: r,
	}
}

func (s *userService) Create(user *models.User) error {
	return errors.Join(
		user.Validate(),
		s.repo.Create(user),
	)
}

// TODO: do a loop to get by provider given
func (s *userService) GetById(id string) (*models.User, error) {
	return s.repo.GetById(id)
}

func (s *userService) AddGenre(userID, genreID string) (*models.Genre, error) {
	user, err := s.GetByGoogleID(userID)
	if err != nil {
		return nil, err
	}

    if len(user.Genres) >= models.UserGenresLimit {
        return nil, models.ErrUserGenresLimitExceeded
    }

	genre, err := s.repo.FirstGenre(genreID)
	if err != nil {
		return nil, err
	}

	if user.HasGenre(genre.Name) {
		return genre, nil
	}

	err = s.repo.AddGenre(user, genre)
	if err != nil {
		return nil, err
	}
	return genre, nil
}

func (s *userService) RemoveGenre(userID, genreID string) (*models.Genre, error) {
	user, err := s.GetByGoogleID(userID)
	if err != nil {
		return nil, err
	}

    fmt.Println(user.Username)

	genre, err := s.repo.FirstGenre(genreID)
	if err != nil {
		return nil, err
	}

	return genre, s.repo.RemoveGenre(user, genre)
}

func (s *userService) GetByGoogleID(googleID string) (*models.User, error) {
	return s.repo.GetByGoogleID(googleID)
}

func (s *userService) GetByEmail(email string) (*models.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) Update(user *models.User) error {
	return s.repo.Update(user)
}

func (s *userService) Delete(userID string) error {
	return s.repo.Delete(userID)
}

func (s *userService) GetAllGenres() ([]*models.Genre, error) {
	return s.repo.GetAllGenres()
}
