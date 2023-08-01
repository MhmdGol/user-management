package service

import "user-management/internal/model"

type UserService interface {
	Create(model.User, model.JwtToken) error
	All() ([]model.User, error)
	ReadByUsername(model.User) (model.User, error)
	UpdateByID(model.User, model.JwtToken) error
	DeleteByID(model.ID, model.JwtToken) error
}

type AuthService interface {
	Login(model.LoginRequest) (model.JwtToken, error)
	Role(t model.JwtToken) (model.Role, error)
	UpdatePassword(model.UpdatePassword) error
}
