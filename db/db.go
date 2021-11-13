package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Should be encapsulated in env variable retrieval
func New() *gorm.DB {
	dsn := "host=localhost user=gorm password=supersecretpassword dbname=deltadb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
