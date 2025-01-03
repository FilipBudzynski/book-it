package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetById(id string) (*models.User, error) {
	user := &models.User{}
	return user, r.db.First(user, id).Error
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	return user, r.db.First(user, "google_id = ?", googleID).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	return user, r.db.First(user, "email = ?", email).Error
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	return users, r.db.Find(&users).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(user models.User) error {
	return r.db.Delete(&user).Error
}
