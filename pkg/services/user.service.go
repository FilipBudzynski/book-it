package services

import (
	"log"

	"github.com/FilipBudzynski/book_it/pkg/entities"
	"gorm.io/gorm"
)

// User provides actions for manipulating users in database.
type User interface {
	Create(u *entities.User) error
	Update(u *entities.User) error
	GetById(id uint) (*entities.User, error)
	GetAll() ([]entities.User, error)
	Delete(u entities.User) error
}

// userService implements the UserService
type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *userService {
	// TODO: user atlas as migration
	if err := db.AutoMigrate(&entities.User{}); err != nil {
		log.Printf("error migrating User entity, err %v", err)
		return nil
	}

	return &userService{
		db: db,
	}
}

func (u *userService) Create(user *entities.User) error {
	return u.db.Create(user).Error
}

func (u *userService) GetById(id uint) (*entities.User, error) {
	var user entities.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userService) GetAll() ([]entities.User, error) {
	var users []entities.User
	if err := u.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *userService) Update(user *entities.User) error {
	return u.db.Save(user).Error
}

func (u *userService) Delete(user entities.User) error {
	return u.db.Delete(user).Error
}
