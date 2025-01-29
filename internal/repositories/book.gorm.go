package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *bookRepository {
	return &bookRepository{
		db: db,
	}
}

func (r *bookRepository) Create(book *models.Book) error {
	var genres []models.Genre

	for _, genreName := range book.Genres {
		var genre models.Genre
		err := r.db.Where("name = ?", genreName.Name).First(&genre).Error
		if err != nil {
			if err.Error() == "record not found" {
				genre = models.Genre{Name: genreName.Name}
				if err := r.db.Create(&genre).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		genres = append(genres, genre)
	}

	book.Genres = genres
	return r.db.Debug().Create(book).Error
}

func (r *bookRepository) Get(id string) (*models.Book, error) {
	book := &models.Book{}
	if err := r.db.First(&book, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return book, nil
}

func (r *bookRepository) Delete(bookID string) error {
	return r.db.Delete(&models.Book{}, bookID).Error
}

func (r *bookRepository) GetByGenre(genre string) ([]*models.Book, error) {
	var books []*models.Book
	err := r.db.Debug().
		Joins("JOIN book_genres ON book_genres.book_id = books.id").
		Joins("JOIN genres ON genres.id = book_genres.genre_id").
		Where("genres.name = ?", genre).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}
