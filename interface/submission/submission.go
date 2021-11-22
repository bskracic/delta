package submission

import "github.com/bSkracic/delta-rest/model"

type Store interface {
	GetSubmissions() ([]model.Submission, error)
	GetSubmission(id uint) (*model.Submission, error)
	CreateSubmission(s *model.Submission) error
	UpdateSubmission(s *model.Submission) error
	DeleteSubmission(s *model.Submission) error
}
