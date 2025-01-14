package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type userBookRepository struct {
	db *gorm.DB
}

func NewUserBookRepository(db *gorm.DB) *userBookRepository {
	return &userBookRepository{
		db: db,
	}
}

func (r *userBookRepository) Create(userBook *models.UserBook) error {
	return r.db.Create(userBook).Error
}

func (r *userBookRepository) Get(id string) (*models.UserBook, error) {
	userBook := &models.UserBook{}
	return userBook, r.db.Preload("Book").First(&userBook, id).Error
}

func (r *userBookRepository) GetAllUserBooks(userId string) ([]*models.UserBook, error) {
	userBooks := []*models.UserBook{}
	return userBooks, r.db.Preload("Book").Preload("ReadingProgress").
		Where("user_google_id = ?", userId).
		Where("deleted_at IS NULL").
		Find(&userBooks).Error
}

func (r *userBookRepository) Update(userBook *models.UserBook) error {
	return r.db.Save(userBook).Error
}

func (r *userBookRepository) Delete(id string) error {
	return r.db.Delete(&models.UserBook{}, id).Error
}

func (r *userBookRepository) DeleteWhereBookId(bookId string) error {
	return r.db.Where("book_id = ?", bookId).Delete(&models.UserBook{}).Error
}

func (r *userBookRepository) Search(userId, query string) ([]*models.UserBook, error) {
	var userBooks []*models.UserBook

	err := r.db.Preload("Book").
		Preload("ReadingProgress").
		Joins("JOIN books ON books.id = user_books.book_id"). // Join the books table
		Where("user_books.user_google_id = ?", userId).       // Filter by current user
		Where("books.title LIKE ?", "%"+query+"%").           // Use the books table for the title
		Where("user_books.deleted_at IS NULL").
		Order("user_books.created_at DESC").
		Find(&userBooks).Error
	if err != nil {
		return nil, err
	}

	return userBooks, nil
}
