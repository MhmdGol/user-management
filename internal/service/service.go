package service

import "user-management/internal/model"

type UserService interface {
	Create(model.User) error
}
