package integration_tests

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/utils"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSession(t *testing.T) {
	t.Helper()
	gothic.Store = sessions.NewCookieStore([]byte("secretstring"))
}

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

func setSession(t *testing.T, req *http.Request, userSession utils.UserSession) {
	t.Helper()
	rec := httptest.NewRecorder()
	err := utils.SetUserSession(rec, req, userSession)
	if err != nil {
		t.Fatalf("Failed to set session: %v", err)
	}
}
