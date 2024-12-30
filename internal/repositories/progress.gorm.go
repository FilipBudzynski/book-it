package repositories

import (
	"github.com/FilipBudzynski/book_it/internal/models"
	"gorm.io/gorm"
)

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) *progressRepository {
	return &progressRepository{
		db: db,
	}
}

func (r *progressRepository) Create(progress models.ReadingProgress) error {
	return r.db.Create(&progress).Error
}

func (r *progressRepository) get(query, data string, preloads ...string) ([]*models.ReadingProgress, error) {
	readingProgress := []*models.ReadingProgress{}
	tx := &gorm.DB{}
	for _, v := range preloads {
		tx = r.db.Preload(v)
	}
	return readingProgress, tx.Find(&readingProgress, query, data).Error
}

func (r *progressRepository) GetById(id string, preloads ...string) (*models.ReadingProgress, error) {
	progress, err := r.get("id = ?", id, preloads...)
	return progress[0], err
}

func (r *progressRepository) GetByUserBookId(userBookId string, preloads ...string) (*models.ReadingProgress, error) {
	progress, err := r.get("user_book_id = ?", userBookId, preloads...)
	return progress[0], err
}

func (r *progressRepository) GetLogById(id string) (*models.DailyProgressLog, error) {
	log := &models.DailyProgressLog{}
	return log, r.db.First(log, id).Error
}

func (r *progressRepository) UpdateLog(log *models.DailyProgressLog) error {
	return r.db.Save(log).Error
}

func (r *progressRepository) Delete(id string) error {
	return r.db.Delete(&models.ReadingProgress{}, id).Error
}
