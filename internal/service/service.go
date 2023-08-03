package service

import (
	"context"
	"user-management/internal/model"
)

type UserService interface {
	Create(context.Context, model.UserInfo) error
	All(context.Context) ([]model.UserInfo, error)
	ReadByUsername(context.Context, model.Username) (model.UserInfo, error)
	UpdateByID(context.Context, model.ID, model.UserInfo) error
	DeleteByID(context.Context, model.ID) error
}

type AuthService interface {
	Login(context.Context, model.Username, model.Password) (model.JwtToken, error)
	Role(model.JwtToken) (model.Role, error)
	Username(model.JwtToken) (model.Username, error)
	UpdatePassword(context.Context, model.UpdatePassword) error
}
