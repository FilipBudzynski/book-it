package services

import (
	"fmt"

	"github.com/FilipBudzynski/book_it/pkg/models"
	"gorm.io/gorm"
)


// userService implements the UserService
type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *userService {
	return &userService{
		db: db,
	}
}

func (u *userService) Create(user *models.User) error {
	return u.db.Create(user).Error
}

// TODO: do a loop to get by providerID or something

func (u *userService) GetById(id string) (*models.User, error) {
	var user models.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *userService) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := u.db.First(&user, "google_id = ?", googleID).Error
	if err == nil {
		return &user, nil
	}

	return nil, err
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

func (u *userService) AddBook(userID, bookID string) error {
	var user models.User
	err := u.db.Preload("Books").First(&user, "google_id = ?", userID).Error
	if err != nil {
		return err
	}

	var userBook models.UserBook
	err = u.db.First(&userBook, "book_id = ?", bookID).Error
	if err == nil {
		return fmt.Errorf("User: %s already has book with id: %s, err: %v", user.Username, bookID, err)
	}

	newUserBook := models.UserBook{
		UserGoogleId: user.GoogleId,
		BookId:       bookID,
		Status:       models.BookStatusNotStarted,
	}

	user.Books = append(user.Books, newUserBook)
	err = u.db.Save(&user).Error
	if err != nil {
		return err
	}

	return nil
}
