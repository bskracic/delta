package execentry

import "github.com/bSkracic/delta-rest/model"

type Store interface {
	GetExecEntries() ([]model.ExecEntry, error)
	GetExecEntry(id uint) (*model.ExecEntry, error)
	GetExecEntriesForSubmission(submissionId uint) ([]model.ExecEntry, error)
	Create(e *model.ExecEntry) error
	Update(e *model.ExecEntry) error
	Delete(e *model.ExecEntry) error
}
