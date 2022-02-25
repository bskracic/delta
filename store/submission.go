package store

import (
	"github.com/bSkracic/delta-rest/model"
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
	return s, ss.db.Preload("User").Preload("Language").First(s).Error
}

func (ss *SubmissionStore) GetSubmission(id uint) (*model.Submission, error) {
	s := model.Submission{}

	if err := ss.db.Preload("User").Preload("Language").First(&s, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (ss *SubmissionStore) GetSubmissionsForAuthor(authorId uint) ([]model.Submission, error) {
	var s []model.Submission

	if err := ss.db.Where(&model.Submission{UserID: authorId}).Preload("User").Preload("Language").Find(&s).Error; err != nil {
		return nil, err
	}
	return s, nil
}

func (ss *SubmissionStore) CreateSubmission(s *model.Submission) error {
	return ss.db.Create(s).Preload("User").Preload("Language").Find(s, s.ID).Error
}

func (ss *SubmissionStore) UpdateSubmission(s *model.Submission) error {
	return ss.db.Model(s).Updates(*s).Error
}

func (ss *SubmissionStore) DeleteSubmission(s *model.Submission) error {
	return ss.db.Delete(s).Error
}
