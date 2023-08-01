package service

import (
	"fmt"
	"io/ioutil"
	"time"
	"user-management/internal/config"
	"user-management/internal/model"
	"user-management/internal/repository"
	"user-management/internal/service"

	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	userRepo      repository.UserRepository
	secretKeyPath string
	publicKeyPath string
}

var _ service.AuthService = (*AuthService)(nil)

func NewAuthService(r repository.UserRepository, conf config.RsaPair) *AuthService {
	return &AuthService{
		userRepo:      r,
		secretKeyPath: conf.SecretKeyPath,
		publicKeyPath: conf.PublicKeyPath,
	}
}

func (as *AuthService) Login(lr model.LoginRequest) (model.JwtToken, error) {
	user, err := as.userRepo.ReadByUsername(model.User{
		Username: lr.Username,
	})
	if err != nil {
		return model.JwtToken{}, err
	}

	if user.Password != lr.Password {
		return model.JwtToken{}, fmt.Errorf("unauthorized")
	}

	expirationTime := time.Now().Add(time.Hour)
	claims := jwt.MapClaims{
		"role": user.Role,
		"exp":  expirationTime.Unix(),
	}

	privateKeyBytes, err := ioutil.ReadFile(as.secretKeyPath)
	if err != nil {
		panic(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return model.JwtToken{}, err
	}

	return model.JwtToken{
		Token: tokenString,
	}, err
}

func (as *AuthService) Role(t model.JwtToken) (model.Role, error) {
	publicKeyBytes, err := ioutil.ReadFile(as.publicKeyPath)
	if err != nil {
		panic(err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		panic(err)
	}

	token, err := jwt.Parse(t.Token, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return model.Role(""), err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return model.Role(""), fmt.Errorf("token has expired")
		}

		role, ok := claims["role"].(string)
		if !ok {
			return model.Role(""), fmt.Errorf("invalid token: role claim not found")
		}

		return model.Role(role), nil
	}

	return model.Role(""), fmt.Errorf("invalid token")
}

func (as *AuthService) UpdatePassword(up model.UpdatePassword) error {
	u, err := as.userRepo.ReadByUsername(model.User{
		Username: up.Username,
	})
	if err != nil {
		return err
	}
	if u.Password != up.OldPassword {
		return fmt.Errorf("not allowed")
	}

	return as.userRepo.UpdateByUsername(model.User{
		Username: up.Username,
		Password: up.NewPassword,
		Version:  u.Version,
	})
}
