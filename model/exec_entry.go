package model

import "github.com/jinzhu/gorm"

type ExecEntry struct {
	gorm.Model
	SubmissionId uint
	Submission   Submission
	Status       string
}
