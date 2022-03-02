package db

import (
	"time"

	"github.com/bSkracic/delta-rest/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	c := config.GetFromEnv()
	for {
		db, err := gorm.Open(postgres.Open(c.ConnString()), &gorm.Config{})

		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
	return db
	}

}
