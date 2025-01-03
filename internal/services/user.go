package services

import (
	"errors"

	"github.com/FilipBudzynski/book_it/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	GetById(id string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(user models.User) error
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

func (s *userService) Delete(user models.User) error {
	return s.repo.Delete(user)
}
