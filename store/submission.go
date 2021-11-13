package store

import (
	"github.com/bSkracic/delta-cli/model"
	"gorm.io/gorm"
)

type SubmissionStore struct {
	db *gorm.DB
}

func NewSubmissionStore(db *gorm.DB) *SubmissionStore {
	return &SubmissionStore{db: db}
}

func (ss *SubmissionStore) GetSubmissions() ([]model.Submission, error) {
	var s []model.Submission
	return s, ss.db.First(s).Error
}

func (ss *SubmissionStore) GetSubmission(id uint) (*model.Submission, error) {
	s := model.Submission{}
	if err := ss.db.First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (ss *SubmissionStore) CreateSubmission(s *model.Submission) error {
	return ss.db.Create(s).Error
}

func (ss *SubmissionStore) UpdateSubmission(s *model.Submission) error {
	return ss.db.Model(s).Updates(*s).Error
}

func (ss *SubmissionStore) DeleteSubmission(s *model.Submission) error {
	return ss.db.Delete(s).Error
}
