package db

import (
	"time"

	"github.com/bSkracic/delta-rest/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	c := config.GetFromEnv()
	//"host=localhost user=gorm password=supersecretpassword dbname=deltadb port=5432 sslmode=disable"
	for {
		db, err := gorm.Open(postgres.Open(c.ConnString()), &gorm.Config{})

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		return db
	}

}
