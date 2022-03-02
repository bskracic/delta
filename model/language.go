package model

import "gorm.io/gorm"

type Language struct {
	gorm.Model
	Name         string
	CompileCmd   string
	ExecuteCmd   string
	MainFileName string
	Image        string
}
