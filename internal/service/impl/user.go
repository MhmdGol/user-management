package service

import (
	"context"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"
	"user-management/pkg"
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

func (us *UserService) Create(ctx context.Context, u model.UserInfo) error {
	hPass, err := pkg.HashPassword(string(u.Password))
	if err != nil {
		return err
	}
	return us.userRepo.Create(ctx, model.User{
		Username: u.Username,
		Password: model.HashedPass(hPass),
		Role:     u.Role,
		City:     u.City,
	})
}

func (us *UserService) All(ctx context.Context) ([]model.UserInfo, error) {
	users, err := us.userRepo.All(ctx)
	if err != nil {
		return []model.UserInfo{}, nil
	}

	uInfos := make([]model.UserInfo, len(users))
	for i, u := range users {
		uInfos[i] = model.UserInfo{
			ID:             u.ID,
			Username:       u.Username,
			Role:           u.Role,
			TimeOfCreation: u.TimeOfCreation,
			City:           u.City,
			Version:        u.Version,
		}
	}

	return uInfos, nil
}

func (us *UserService) ReadByUsername(ctx context.Context, u model.Username) (model.UserInfo, error) {
	user, err := us.userRepo.ReadByUsername(ctx, u)
	if err != nil {
		return model.UserInfo{}, nil
	}

	return model.UserInfo{
		ID:             user.ID,
		Username:       user.Username,
		Role:           user.Role,
		TimeOfCreation: user.TimeOfCreation,
		City:           user.City,
		Version:        user.Version,
	}, nil
}

func (us *UserService) UpdateByID(ctx context.Context, id model.ID, u model.UserInfo) error {
	hPass, err := pkg.HashPassword(string(u.Password))
	if err != nil {
		return err
	}

	return us.userRepo.UpdateByID(ctx, model.User{
		ID:       id,
		Username: u.Username,
		Password: model.HashedPass(hPass),
		Role:     u.Role,
		City:     u.City,
		Version:  u.Version,
	})
}

func (us *UserService) DeleteByID(ctx context.Context, id model.ID) error {
	return us.userRepo.DeleteByID(ctx, id)
}
