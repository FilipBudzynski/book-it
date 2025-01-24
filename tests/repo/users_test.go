package test_repo

import (
	"log"
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	// Open an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Enable foreign key support (important for SQLite in-memory DB)
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		log.Fatalf("failed to enable foreign key support: %v", err)
	}

	// Run migrations to create the schema for your models
	err = db.AutoMigrate(models.MigrateModels...)
	if err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	// Return the database connection
	return db
}

func TestCascadeDeleteUser(t *testing.T) {
	// Setup
	db := setupTestDB()

	// Create some books first
	book1 := models.Book{ID: "book1"}
	book2 := models.Book{ID: "book2"}

	// Insert the books into the database
	if err := db.Create(&book1).Error; err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if err := db.Create(&book2).Error; err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}

	// Create a user with associated UserBooks
	//
	user := models.User{
		GoogleId: "123",
		Username: "testuser",
		Email:    "test@example.com",
		Books: []models.UserBook{
			{
				BookID:       "book1", // Now the BookID refers to an existing Book
				UserGoogleId: "123",
			},
			{
				BookID:       "book2", // Now the BookID refers to an existing Book
				UserGoogleId: "123",
			},
		},
	}

	user2 := models.User{GoogleId: "321"}

	// Create the user
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Check if UserBooks were created
	var count int64
	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 UserBooks, found %d", count)
	}

	repo := repositories.NewUserRepository(db)
	// Delete the user
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	db.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 User, found %d", count)
	}

	// Check if UserBooks were deleted
	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 UserBooks, found %d", count)
	}
}

func TestCascadeDeleteExchangeRequest(t *testing.T) {
	// Setup
	//
	db := setupTestDB()

	// Create some books
	book1 := models.Book{ID: "book1"}
	book2 := models.Book{ID: "book2"}

	if err := db.Create(&book1).Error; err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if err := db.Create(&book2).Error; err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}

	// Create a user
	user := models.User{
		GoogleId: "123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Create an exchange request associated with the user and books
	exchangeRequest := models.ExchangeRequest{
		UserGoogleId:  "123",
		UserEmail:     "test@example.com",
		DesiredBookID: "book1",
	}

	if err := db.Create(&exchangeRequest).Error; err != nil {
		t.Fatalf("Failed to create exchange request: %v", err)
	}

	// Before deleting the exchange request, check that it exists and is linked to the user
	var count int64
	db.Model(&models.ExchangeRequest{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 ExchangeRequest, found %d", count)
	}

	repo := repositories.NewUserRepository(db)
	// Delete the user
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	db.Model(&models.User{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 User, found %d", count)
	}

	// Check that the ExchangeRequest was deleted
	db.Model(&models.ExchangeRequest{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 0 {
		t.Fatalf("Expected 0 ExchangeRequests, found %d", count)
	}
}
