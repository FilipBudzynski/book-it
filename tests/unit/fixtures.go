package unit

import (
	"log"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:?_fk=1"), &gorm.Config{})
	require.NoError(t, err, "Failed to connect to in-memory SQLite")

	if err := db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Fatalf("failed to enable foreign key support: %v", err)
	}

	err = db.AutoMigrate(models.MigrateModels...)
	require.NoError(t, err, "Failed to migrate database")

	return db, func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
