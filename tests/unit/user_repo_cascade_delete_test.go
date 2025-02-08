package unit

import (
	"testing"

	"github.com/FilipBudzynski/book_it/internal/models"
	"github.com/FilipBudzynski/book_it/internal/repositories"
)

func TestCascadeDeleteUserBooks(t *testing.T) {
	db, cleanUp := setupTestDB(t)
	defer cleanUp()

	book1 := models.Book{ID: "book1"}
	book2 := models.Book{ID: "book2"}

	if err := db.Create(&book1).Error; err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if err := db.Create(&book2).Error; err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}

	user := models.User{
		GoogleId: "123",
		Username: "testuser",
		Email:    "test@example.com",
		Books: []models.UserBook{
			{
				BookID:       "book1",
				UserGoogleId: "123",
			},
			{
				BookID:       "book2",
				UserGoogleId: "123",
			},
		},
	}

	user2 := models.User{
		GoogleId: "321",
		Email:    "test2@example.com",
		Username: "testuser2",
		Books: []models.UserBook{
			{
				BookID:       "book1",
				UserGoogleId: "321",
			},
			{
				BookID:       "book2",
				UserGoogleId: "321",
			},
		},
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	var count int64
	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 UserBooks, found %d", count)
	}

	repo := repositories.NewUserRepository(db)
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}
	db.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 User, found %d", count)
	}

	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 UserBooks, found %d", count)
	}
}

func TestCascadeDeleteExchangeRequest(t *testing.T) {
	db, cleanUp := setupTestDB(t)
	defer cleanUp()

	book1 := models.Book{ID: "book1"}
	book2 := models.Book{ID: "book2"}

	if err := db.Create(&book1).Error; err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if err := db.Create(&book2).Error; err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}

	user := models.User{
		GoogleId: "123",
		Username: "testuser",
		Email:    "test@example.com",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	exchangeRequest := models.ExchangeRequest{
		UserGoogleId:  "123",
		UserEmail:     "test@example.com",
		DesiredBookID: "book1",
		OfferedBooks: []models.OfferedBook{
			{
				BookID: "book2",
			},
		},
	}

	if err := db.Create(&exchangeRequest).Error; err != nil {
		t.Fatalf("Failed to create exchange request: %v", err)
	}

	var count int64
	db.Model(&models.ExchangeRequest{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 1 {
		t.Fatalf("Expected 1 ExchangeRequest, found %d", count)
	}

	repo := repositories.NewUserRepository(db)
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	db.Model(&models.User{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 User, found %d", count)
	}

	db.Model(&models.ExchangeRequest{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 0 {
		t.Fatalf("Expected 0 ExchangeRequests, found %d", count)
	}
}

func TestCascadeDeleteUserBooksWithProgress(t *testing.T) {
	db, cleanUp := setupTestDB(t)
	defer cleanUp()

	book1 := models.Book{ID: "book1"}
	book2 := models.Book{ID: "book2"}

	if err := db.Create(&book1).Error; err != nil {
		t.Fatalf("Failed to create book1: %v", err)
	}
	if err := db.Create(&book2).Error; err != nil {
		t.Fatalf("Failed to create book2: %v", err)
	}

	user := models.User{GoogleId: "123", Username: "testuser", Email: "test@example.com"}

	user2 := models.User{GoogleId: "321", Email: "test2@example.com", Username: "testuser2"}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if err := db.Create(&user2).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	user_book_1 := models.UserBook{BookID: "book1", UserGoogleId: "123"}

	user_book_2 := models.UserBook{BookID: "book2", UserGoogleId: "123"}

	if err := db.Create(&user_book_1).Error; err != nil {
		t.Fatalf("Failed to create user book: %v", err)
	}
	if err := db.Create(&user_book_2).Error; err != nil {
		t.Fatalf("Failed to create user book: %v", err)
	}

	progress_1 := models.ReadingProgress{
		UserBookID: user_book_1.ID,
		TotalPages: 400,
		DailyProgress: []models.DailyProgressLog{
			{
				PagesRead: 100,
			},
			{
				PagesRead: 200,
			},
		},
	}
	if err := db.Create(&progress_1).Error; err != nil {
		t.Fatalf("Failed to create progress: %v", err)
	}

	var count int64
	db.Model(&models.ReadingProgress{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 Progress, found %d", count)
	}

	count = 0
	db.Model(&models.DailyProgressLog{}).Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 reading progresses, found %d", count)
	}

	count = 0
	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 2 {
		t.Errorf("Expected 2 UserBooks, found %d", count)
	}

	repo := repositories.NewUserRepository(db)
	err := repo.Delete("123")
	if err != nil {
		t.Fatalf("Failed to delete user: %v", err)
	}

	count = 0
	db.Model(&models.User{}).Count(&count)
	if count != 1 {
		t.Errorf("Expected 1 User, found %d", count)
	}

	count = 1
	db.Model(&models.UserBook{}).Where("user_google_id = ?", "123").Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 UserBooks, found %d", count)
	}

	count = 1
	db.Model(&models.ReadingProgress{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 reading progresses, found %d", count)
	}

	count = 1
	db.Model(&models.DailyProgressLog{}).Count(&count)
	if count != 0 {
		t.Errorf("Expected 0 daily logs, found %d", count)
	}
}
