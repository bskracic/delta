package store

import (
	"github.com/bSkracic/delta-rest/model"
	"gorm.io/gorm"
)

type ExecEntryStore struct {
	db *gorm.DB
}

func NewExecEntryStore(db *gorm.DB) *ExecEntryStore {
	return &ExecEntryStore{
		db: db,
	}
}

func (es *ExecEntryStore) GetExecEntries() ([]model.ExecEntry, error) {
	var e []model.ExecEntry
	return e, es.db.Preload("Submission").Preload("Submission.Author").Preload("Submission.Language").First(&e).Error
}

func (es *ExecEntryStore) GetExecEntry(id uint) (*model.ExecEntry, error) {
	e := model.ExecEntry{}
	// .Preload("Submission.User").Preload("Submission.Language")
	if err := es.db.Preload("Submission.User").Preload("Submission.Language").First(&e, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, gorm.ErrDryRunModeUnsupported
	}

	// var s model.Submission
	// if err := es.db.Preload("Author").Preload("Language").First(&s, e.SubmissionId).Error; err != nil {
	// 	return nil, err
	// }

	// e.Submission = s
	return &e, nil
}

func (es *ExecEntryStore) Create(e *model.ExecEntry) error {
	return es.db.Create(e).Error
}

func (es *ExecEntryStore) Update(e *model.ExecEntry) error {
	return es.db.Model(e).Updates(*e).Error
}

func (es *ExecEntryStore) Delete(e *model.ExecEntry) error {
	return es.db.Delete(e).Error
}
