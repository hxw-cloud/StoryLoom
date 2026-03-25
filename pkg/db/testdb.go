package db

import (
	"log"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitTestDB initializes an in-memory SQLite database specifically for unit testing.
// It bypasses the need for file cleanup and ensures tests run isolated and fast.
//
// Why this approach? TDD requires fast feedback loops. File-based databases
// introduce I/O latency and potential state leakage between tests. By using
// "file::memory:?cache=shared", we get a blazing fast, ephemeral database.
func InitTestDB(models ...interface{}) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Silent, // Keep test output clean
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Use in-memory SQLite for testing
	testDB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	// Auto-migrate the provided domain models to prep the schema
	if len(models) > 0 {
		err = testDB.AutoMigrate(models...)
		if err != nil {
			log.Fatalf("failed to migrate test database: %v", err)
		}
	}

	// Assign to the global DB variable used by the application handlers
	DB = testDB
	return testDB
}
