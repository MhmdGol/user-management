package service

import (
	"fmt"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"
)

type UserService struct {
	userRepo    repository.UserRepository
	authService service.AuthService
}

var _ service.UserService = (*UserService)(nil)

func NewUserService(r repository.UserRepository, a service.AuthService) *UserService {
	return &UserService{
		userRepo:    r,
		authService: a,
	}
}

func (us *UserService) Create(u model.User, t model.JwtToken) error {
	r, err := us.authService.Role(t)
	if err != nil {
		return err
	}
	if r != "admin" {
		return fmt.Errorf("not allowed")
	}

	return us.userRepo.Create(u)
}

func (us *UserService) All() ([]model.User, error) {
	return us.userRepo.All()
}

func (us *UserService) ReadByUsername(u model.User) (model.User, error) {
	return us.userRepo.ReadByUsername(u)
}

func (us *UserService) UpdateByID(u model.User, t model.JwtToken) error {
	r, err := us.authService.Role(t)
	if err != nil {
		return err
	}
	if r != "admin" && r != "staff" {
		return fmt.Errorf("not allowed")
	}

	return us.userRepo.UpdateByID(u)
}

func (us *UserService) DeleteByID(id model.ID, t model.JwtToken) error {
	r, err := us.authService.Role(t)
	if err != nil {
		return err
	}
	if r != "admin" {
		return fmt.Errorf("not allowed")
	}

	return us.userRepo.DeleteByID(id)
}
