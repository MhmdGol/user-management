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

func (s *UserServiceServer) Create(ctx context.Context, req *protobuf.UserRequest) (*protobuf.Response, error) {
	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Role:     model.Role(req.Role),
		City:     req.City,
	}

	token := model.JwtToken{
		Token: req.Token.Token,
	}

	err := s.userSrv.Create(user, token)

	return &protobuf.Response{
		Success: true,
	}, err
}

func (s *UserServiceServer) Login(ctx context.Context, req *protobuf.LoginInfo) (*protobuf.JwtToken, error) {
	login := model.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := s.authSrv.Login(login)

	return &protobuf.JwtToken{
		Token: token.Token,
	}, err
}

func (s *UserServiceServer) All(ctx context.Context, req *protobuf.Empty) (*protobuf.Users, error) {
	users, _ := s.userSrv.All()

	result := make([]*protobuf.User, len(users))
	for i, u := range users {
		result[i] = &protobuf.User{
			Id:             string(u.ID),
			Username:       u.Username,
			Password:       u.Password,
			TimeOfCreation: u.TimeOfCreation.String(),
			Role:           string(u.Role),
			City:           u.City,
			Version:        int32(u.Version),
		}
	}
	return &protobuf.Users{
		Users: result,
	}, nil
}

func (s *UserServiceServer) ReadByUsername(ctx context.Context, req *protobuf.Username) (*protobuf.User, error) {
	u, _ := s.userSrv.ReadByUsername(model.User{Username: req.Username})

	return &protobuf.User{
		Id:             string(u.ID),
		Username:       u.Username,
		Password:       u.Password,
		TimeOfCreation: u.TimeOfCreation.String(),
		Role:           string(u.Role),
		City:           u.City,
		Version:        int32(u.Version),
	}, nil
}

func (s *UserServiceServer) UpdateByID(ctx context.Context, up *protobuf.Update) (*protobuf.Response, error) {
	s.userSrv.UpdateByID(model.User{
		ID:       model.ID(up.User.Id),
		Username: up.User.Username,
		Password: up.User.Password,
		Role:     model.Role(up.User.Role),
		City:     up.User.City,
		Version:  int(up.User.Version),
	}, model.JwtToken{
		Token: up.Token.Token,
	})

	return &protobuf.Response{
		Success: true,
	}, nil
}

func (s *UserServiceServer) DeleteByID(ctx context.Context, de *protobuf.Delete) (*protobuf.Response, error) {
	s.userSrv.DeleteByID(model.ID(de.Id.Id), model.JwtToken{
		Token: de.Token.Token,
	})

	return &protobuf.Response{
		Success: true,
	}, nil
}

func (s *UserServiceServer) UpdatePassword(ctx context.Context, upass *protobuf.UpdatePass) (*protobuf.Response, error) {
	s.authSrv.UpdatePassword(model.UpdatePassword{
		Username: upass.Username,
		OldPassword: upass.Oldpass,
		NewPassword: upass.Newpass,
	})

	return &protobuf.Response{
		Success: true,
	}, nil
}
