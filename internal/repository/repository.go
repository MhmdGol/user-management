package repository

import "user-management/internal/model"

type UserRepository interface {
	Create(model.User) error
	All() ([]model.User, error)
	ReadByUsername(model.User) (model.User, error)
	UpdateByID(model.User) error
	UpdateByUsername(model.User) error
	DeleteByID(model.ID) error
}
