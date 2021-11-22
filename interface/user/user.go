package user

import "github.com/bSkracic/delta-rest/model"

type Store interface {
	GetById(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create(u *model.User) error
	Update(u *model.User) error
	Delete(u *model.User) error
}
