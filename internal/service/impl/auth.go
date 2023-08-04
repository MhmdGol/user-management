package service

import (
	"context"
	"fmt"
	"time"
	"user-management/internal/jwtpkg"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"
	"user-management/pkg"
)

type AuthService struct {
	userRepo repository.UserRepository
	JwtToken *jwtpkg.JwtToken
}

var _ service.AuthService = (*AuthService)(nil)

func NewAuthService(r repository.UserRepository, j *jwtpkg.JwtToken) *AuthService {
	return &AuthService{
		userRepo: r,
		JwtToken: j,
	}
}

func (as *AuthService) Login(ctx context.Context, u model.Username, p model.Password) (model.JwtToken, error) {
	user, err := as.userRepo.ReadByUsername(ctx, u)
	if err != nil {
		return "", err
	}

	err = pkg.ValidatePassword(string(user.Password), string(p))
	if err != nil {
		return "", fmt.Errorf("unauthorized")
	}

	tokenString, err := as.JwtToken.MakeToken(model.TokenClaim{
		Username:       u,
		Role:           user.Role,
		ExpirationTime: time.Now().Add(time.Hour),
	})
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func (as *AuthService) Role(t model.JwtToken) (model.Role, error) {
	claims, err := as.JwtToken.ExtractClaims(t)
	if err != nil {
		return model.Role(""), err
	}

	return claims.Role, nil
}

func (as *AuthService) Username(t model.JwtToken) (model.Username, error) {
	claims, err := as.JwtToken.ExtractClaims(t)
	if err != nil {
		return model.Username(""), err
	}

	return claims.Username, nil
}

func (as *AuthService) UpdatePassword(ctx context.Context, up model.UpdatePassword) error {
	user, err := as.userRepo.ReadByUsername(ctx, up.Username)
	if err != nil {
		return err
	}

	hOldPass, _ := pkg.HashPassword(string(up.OldPassword))
	if user.Password != model.HashedPass(hOldPass) {
		return fmt.Errorf("not allowed")
	}

	hNewPass, _ := pkg.HashPassword(string(up.NewPassword))
	return as.userRepo.UpdateByUsername(ctx, model.User{
		Username: up.Username,
		Password: model.HashedPass(hNewPass),
		Version:  user.Version,
	})
}
