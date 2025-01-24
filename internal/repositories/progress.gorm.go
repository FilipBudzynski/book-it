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
	return r.db.Debug().Create(&progress).Error
}

func (r *progressRepository) GetById(id string) (*models.ReadingProgress, error) {
	progress := &models.ReadingProgress{}
	return progress, r.db.Preload("DailyProgress").First(progress, "id = ?", id).Error
}

func (r *progressRepository) GetByUserBookId(userBookId string) (*models.ReadingProgress, error) {
	progress := &models.ReadingProgress{}
	return progress, r.db.Preload("DailyProgress").First(progress, "user_book_id = ?", userBookId).Error
}

func (r *progressRepository) GetLogById(id string) (*models.DailyProgressLog, error) {
	log := &models.DailyProgressLog{}
	return log, r.db.First(log, id).Error
}

func (r *progressRepository) Update(progress *models.ReadingProgress) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(progress).Error
}

func (r *progressRepository) UpdateLog(log *models.DailyProgressLog) error {
	return r.db.Save(log).Error
}

func (r *progressRepository) Delete(id string) error {
	return r.db.Delete(&models.ReadingProgress{}, id).Error
}
