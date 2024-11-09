package services

import (
	"log"

	"github.com/FilipBudzynski/book_it/pkg/models"
	"gorm.io/gorm"
)

// UserService provides actions for managing Users.
type UserService interface {
	Create(u *models.User) error
	Update(u *models.User) error
	GetById(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Delete(u models.User) error
}

// userService implements the UserService
type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *userService {
	// TODO: user atlas as migration
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Printf("error migrating User entity, err %v", err)
		return nil
	}

	return &userService{
		db: db,
	}
}

func (u *userService) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *userService) GetById(id string) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userService) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, "email = ?", email).Error
	if err == nil {
		return &user, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return nil, err
}

func (u *userService) GetAll() ([]models.User, error) {
	var users []models.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userService) Update(user *models.User) error {
	return u.db.Save(user).Error
}

func (u *userService) Delete(user models.User) error {
	return u.db.Delete(user).Error
}
