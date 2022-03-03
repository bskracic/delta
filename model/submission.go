package model

import "github.com/jinzhu/gorm"

type Submission struct {
	gorm.Model
	MainFile   []byte
	UserID     *uint
	User       User
	LanguageID uint
	Language   Language
}
