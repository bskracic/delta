package store

import (
	"github.com/bSkracic/delta-rest/model"
	"gorm.io/gorm"
)

type LanguageStore struct {
	db *gorm.DB
}

func NewLanguageStore(db *gorm.DB) *LanguageStore {
	return &LanguageStore{
		db: db,
	}
}

func (ls *LanguageStore) GetAll() ([]model.Language, error) {
	var l []model.Language
	return l, ls.db.First(l).Error
}

func (ls *LanguageStore) GetById(id uint) (*model.Language, error) {
	l := model.Language{}
	if err := ls.db.First(&l, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &l, nil
}

func (ls *LanguageStore) Create(l *model.Language) error {
	return ls.db.Create(l).Error
}

func (ls *LanguageStore) Update(l *model.Language) error {
	return ls.db.Model(l).Updates(*l).Error
}

func (ls *LanguageStore) Delete(l *model.Language) error {
	return ls.db.Delete(l).Error
}
