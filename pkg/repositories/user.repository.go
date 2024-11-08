package repositories

import (
	"github.com/FilipBudzynski/book_it/pkg/models"
	"gorm.io/gorm"
)

// UserRepository provides actions for manipulating users in database.
type UserRepository interface {
	Create(u *models.User) error
	GetBy(attribute any) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(u *models.User) error
	Delete(u models.User) error
}

// userRepo implements UserRepository
type userRepo struct {
	db *gorm.DB
}
