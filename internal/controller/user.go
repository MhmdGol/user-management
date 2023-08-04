package controller

import (
	"context"
	"fmt"
	"user-management/internal/model"
	"user-management/internal/proto"
	"user-management/internal/proto/protoconnect"
	service "user-management/internal/service"

	connect_go "github.com/bufbuild/connect-go"
)

type UserServiceServer struct {
	UserSrv service.UserService
	AuthSrv service.AuthService
}

var _ protoconnect.UserServiceHandler = (*UserServiceServer)(nil)

func (c *UserServiceServer) CreateUserService(ctx context.Context, req *connect_go.Request[proto.CreateUserServiceRequest]) (*connect_go.Response[proto.CreateUserServiceResponse], error) {
	token := req.Header().Values("Authorization")[0]

	role, err := c.AuthSrv.Role(model.JwtToken(token))
	if err != nil {
		return &connect_go.Response[proto.CreateUserServiceResponse]{
			Msg: &proto.CreateUserServiceResponse{
				Success: false,
			},
		}, connect_go.NewError(connect_go.Code(16), err)
	}

	if string(role) != "admin" {
		return &connect_go.Response[proto.CreateUserServiceResponse]{
			Msg: &proto.CreateUserServiceResponse{
				Success: false,
			},
		}, connect_go.NewError(connect_go.Code(7), fmt.Errorf("not admin"))
	}

	err = c.UserSrv.Create(ctx, model.UserInfo{
		Username: model.Username(req.Msg.Username),
		Password: model.Password(req.Msg.Password),
		Role:     model.Role(req.Msg.Role),
		City:     req.Msg.City,
	})
	if err != nil {
		return &connect_go.Response[proto.CreateUserServiceResponse]{
			Msg: &proto.CreateUserServiceResponse{
				Success: false,
			},
		}, connect_go.NewError(connect_go.Code(3), err)
	}

	return &connect_go.Response[proto.CreateUserServiceResponse]{
		Msg: &proto.CreateUserServiceResponse{
			Success: false,
		},
	}, nil
}

// GetAllUsersService calls usermanagement.v1.UserService.GetAllUsersService.
func (c *UserServiceServer) GetAllUsersService(ctx context.Context, req *connect_go.Request[proto.GetAllUsersServiceRequest]) (*connect_go.Response[proto.GetAllUsersServiceResponse], error) {
	return nil, nil
}

// GetInfoService calls usermanagement.v1.UserService.GetInfoService.
func (c *UserServiceServer) GetInfoService(ctx context.Context, req *connect_go.Request[proto.GetInfoServiceRequest]) (*connect_go.Response[proto.GetInfoServiceResponse], error) {
	return nil, nil
}

// UpdateByIdService calls usermanagement.v1.UserService.UpdateByIdService.
func (c *UserServiceServer) UpdateByIdService(ctx context.Context, req *connect_go.Request[proto.UpdateByIdServiceRequest]) (*connect_go.Response[proto.UpdateByIdServiceResponse], error) {
	return nil, nil
}

// DeleteByIdService calls usermanagement.v1.UserService.DeleteByIdService.
func (c *UserServiceServer) DeleteByIdService(ctx context.Context, req *connect_go.Request[proto.DeleteByIdServiceRequest]) (*connect_go.Response[proto.DeleteByIdServiceResponse], error) {
	return nil, nil
}

// func (s *UserServiceServer) Create(ctx context.Context, req *protobuf.UserRequest) (*protobuf.Response, error) {
// 	user := model.User{
// 		Username: req.Username,
// 		Password: req.Password,
// 		Role:     model.Role(req.Role),
// 		City:     req.City,
// 	}

// 	token := model.JwtToken{
// 		Token: req.Token.Token,
// 	}

// 	err := s.userSrv.Create(user, token)

// 	return &protobuf.Response{
// 		Success: true,
// 	}, err
// }

// func (s *UserServiceServer) Login(ctx context.Context, req *protobuf.LoginInfo) (*protobuf.JwtToken, error) {
// 	login := model.LoginRequest{
// 		Username: req.Username,
// 		Password: req.Password,
// 	}

// 	token, err := s.authSrv.Login(login)

// 	return &protobuf.JwtToken{
// 		Token: token.Token,
// 	}, err
// }

// func (s *UserServiceServer) All(ctx context.Context, req *protobuf.Empty) (*protobuf.Users, error) {
// 	users, _ := s.userSrv.All()

// 	result := make([]*protobuf.User, len(users))
// 	for i, u := range users {
// 		result[i] = &protobuf.User{
// 			Id:             string(u.ID),
// 			Username:       u.Username,
// 			Password:       u.Password,
// 			TimeOfCreation: u.TimeOfCreation.String(),
// 			Role:           string(u.Role),
// 			City:           u.City,
// 			Version:        int32(u.Version),
// 		}
// 	}
// 	return &protobuf.Users{
// 		Users: result,
// 	}, nil
// }

// func (s *UserServiceServer) ReadByUsername(ctx context.Context, req *protobuf.Username) (*protobuf.User, error) {
// 	u, _ := s.userSrv.ReadByUsername(model.User{Username: req.Username})

// 	return &protobuf.User{
// 		Id:             string(u.ID),
// 		Username:       u.Username,
// 		Password:       u.Password,
// 		TimeOfCreation: u.TimeOfCreation.String(),
// 		Role:           string(u.Role),
// 		City:           u.City,
// 		Version:        int32(u.Version),
// 	}, nil
// }

// func (s *UserServiceServer) UpdateByID(ctx context.Context, up *protobuf.Update) (*protobuf.Response, error) {
// 	s.userSrv.UpdateByID(model.User{
// 		ID:       model.ID(up.User.Id),
// 		Username: up.User.Username,
// 		Password: up.User.Password,
// 		Role:     model.Role(up.User.Role),
// 		City:     up.User.City,
// 		Version:  int(up.User.Version),
// 	}, model.JwtToken{
// 		Token: up.Token.Token,
// 	})

// 	return &protobuf.Response{
// 		Success: true,
// 	}, nil
// }

// func (s *UserServiceServer) DeleteByID(ctx context.Context, de *protobuf.Delete) (*protobuf.Response, error) {
// 	s.userSrv.DeleteByID(model.ID(de.Id.Id), model.JwtToken{
// 		Token: de.Token.Token,
// 	})

// 	return &protobuf.Response{
// 		Success: true,
// 	}, nil
// }

// func (s *UserServiceServer) UpdatePassword(ctx context.Context, upass *protobuf.UpdatePass) (*protobuf.Response, error) {
// 	s.authSrv.UpdatePassword(model.UpdatePassword{
// 		Username:    upass.Username,
// 		OldPassword: upass.Oldpass,
// 		NewPassword: upass.Newpass,
// 	})

// 	return &protobuf.Response{
// 		Success: true,
// 	}, nil
// }
