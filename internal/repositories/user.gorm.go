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

func (r *userRepository) getPreloads() *gorm.DB {
	return r.db.Preload("Genres").Preload("Location")
}

func (r *userRepository) GetById(id string) (*models.User, error) {
	user := &models.User{}
	return user, r.getPreloads().First(user, id).Error
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	return user, r.getPreloads().First(user, "google_id = ?", googleID).Error
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	return user, r.getPreloads().First(user, "email = ?", email).Error
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	return users, r.db.Find(&users).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error
}

func (r *userRepository) Delete(userID string) error {
	return r.db.Debug().Unscoped().Where("google_id = ?", userID).Delete(&models.User{}).Error
}

func (r *userRepository) AddGenre(user *models.User, genre *models.Genre) error {
	return r.db.Model(user).Association("Genres").Append(genre)
}

func (r *userRepository) RemoveGenre(user *models.User, genre *models.Genre) error {
	return r.db.Debug().Model(user).Association("Genres").Delete(genre)

// languages := user.Languages 
// DB.Model(&user).Association("Languages").Clear()
// user.Languages = languages
}

func (r *userRepository) FindOrCreateGenre(genreName string) (*models.Genre, error) {
	genre := &models.Genre{}
	return genre, r.db.FirstOrCreate(genre, models.Genre{Name: genreName}).Error
}

func (r *userRepository) FirstGenre(genreID string) (*models.Genre, error) {
	genre := &models.Genre{}
	return genre, r.db.First(genre, "id = ?", genreID).Error
}

func (r *userRepository) GetAllGenres() ([]*models.Genre, error) {
	var genres []*models.Genre
	return genres, r.db.Find(&genres).Error
}
