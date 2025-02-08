package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/FilipBudzynski/book_it/internal/models"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const openConnection = 40

var (
	dburl      = os.Getenv("DB_URL")
	dbInstance *gorm.DB
)

func init() {
	db, err := gorm.Open(sqlite.Open(dburl+"?_fk=1"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if err := db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Fatalf("failed to enable foreign key support: %v", err)
	}

	err = db.AutoMigrate(models.MigrateModels...)
	if err != nil {
		panic("failed to migrate database")
	}
}

type Service interface {
	Health() map[string]string
	Close() error
}

type Repository struct {
	Db *gorm.DB
}

func New() *gorm.DB {
	if dbInstance != nil {
		return dbInstance
	}

	db, err := gorm.Open(sqlite.Open(dburl+"?_fk=1"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Exec("PRAGMA foreign_keys = ON", nil).Error; err != nil {
		log.Fatalf("failed to enable foreign key support: %v", err)
	}

	return db
}

func (s *Repository) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	db, err := s.Db.DB()
	if err != nil {
		log.Fatalf("gorm conversion to db failed, err: %v", err)
	}
	err = db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > openConnection {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

func (s *Repository) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	db, err := s.Db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
