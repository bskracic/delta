package model

import "github.com/jinzhu/gorm"

type Submission struct {
	gorm.Model
	Language string
	MainFile []byte
}
