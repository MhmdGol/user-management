package service

import (
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"
)

type UserService struct {
	userRepo repository.UserRepository
}

var _ service.UserService = (*UserService)(nil)

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{
		userRepo: r,
	}
}

func (us *UserService) Create(u model.User) error {
	return us.userRepo.Create(u)
}
