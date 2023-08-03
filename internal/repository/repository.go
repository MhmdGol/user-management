package repository

import (
	"context"
	"user-management/internal/model"
)

type UserRepository interface {
	Create(context.Context, model.User) error
	All(context.Context) ([]model.User, error)
	ReadByUsername(context.Context, model.Username) (model.User, error)
	UpdateByID(context.Context, model.User) error
	UpdateByUsername(context.Context, model.User) error
	DeleteByID(context.Context, model.ID) error
}
