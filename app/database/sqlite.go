package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (db *gorm.DB, close func() error) {
	db, err := gorm.Open(sqlite.Open("klines.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %s", err)
	}

	return db, sqlDB.Close
}
