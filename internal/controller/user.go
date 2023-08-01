package controller

import (
	"context"
	"user-management/internal/model"
	protobuf "user-management/internal/protobuf/user"
	service "user-management/internal/service"
)

type UserServiceServer struct {
	srv service.UserService
	protobuf.UnimplementedUserServiceServer
}

func NewUserServiceServer(s service.UserService) *UserServiceServer {
	return &UserServiceServer{
		srv: s,
	}
}

var _ protobuf.UserServiceServer = (*UserServiceServer)(nil)

func (s *UserServiceServer) Create(ctx context.Context, req *protobuf.UserRequest) (*protobuf.UserResponse, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     req.Role,
		City:     req.City,
	}

	err := s.srv.Create(user)

	return &protobuf.UserResponse{
		Success: true,
	}, err
}
