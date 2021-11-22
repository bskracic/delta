package language

import "github.com/bSkracic/delta-rest/model"

type Store interface {
	GetAll() ([]model.Language, error)
	GetById(id uint) (*model.Language, error)
	Create(l *model.Language) error
	Update(l *model.Language) error
	Delete(l *model.Language) error
}
