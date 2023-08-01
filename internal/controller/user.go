package controller

import (
	"context"
	"user-management/internal/model"
	protobuf "user-management/internal/protobuf/user"
	service "user-management/internal/service"
)

type UserServiceServer struct {
	userSrv service.UserService
	authSrv service.AuthService
	protobuf.UnimplementedUserServiceServer
}

func NewUserServiceServer(us service.UserService, as service.AuthService) *UserServiceServer {
	return &UserServiceServer{
		userSrv: us,
		authSrv: as,
	}
}

var _ protobuf.UserServiceServer = (*UserServiceServer)(nil)

func (s *UserServiceServer) Create(ctx context.Context, req *protobuf.UserRequest) (*protobuf.UserResponse, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     model.Role(req.Role),
		City:     req.City,
	}

	err := s.userSrv.Create(user, model.JwtToken{
		Token: "",
	})

	return &protobuf.UserResponse{
		Success: true,
	}, err
}

// func (s *UserServiceServer) Login(ctx context.Context, req *protobuf.LoginInfo) (*protobuf.JwtTokenResponse, error) {
// 	login := model.LoginRequest{
// 		Username: req.Username,
// 		Password: req.Password,
// 	}

// 	token, err :=

// }
