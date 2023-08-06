package controller

import (
	"context"
	"user-management/internal/model"
	userapiv1 "user-management/internal/proto/usermgt/userapi/v1"
	service "user-management/internal/service"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserServiceServer struct {
	userSrv service.UserService
	authSrv service.AuthService
	userapiv1.UnimplementedUserServiceServer
}

func NewUserServiceServer(us service.UserService, as service.AuthService) *UserServiceServer {
	return &UserServiceServer{
		userSrv: us,
		authSrv: as,
	}
}

var _ userapiv1.UserServiceServer = (*UserServiceServer)(nil)

func (s *UserServiceServer) CreateUserService(ctx context.Context, req *userapiv1.CreateUserServiceRequest) (*userapiv1.CreateUserServiceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &userapiv1.CreateUserServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	token, found := md["authorization"]
	if !found {
		return &userapiv1.CreateUserServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	role, err := s.authSrv.Role(model.JwtToken(token[0]))
	if err != nil {
		return &userapiv1.CreateUserServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	if role != "admin" {
		return &userapiv1.CreateUserServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "not allowed")
	}

	user := model.UserInfo{
		Username: model.Username(req.Username),
		Password: model.Password(req.Password),
		Role:     model.Role(req.Role),
		City:     req.City,
	}

	err = s.userSrv.Create(ctx, user)
	if err != nil {
		return &userapiv1.CreateUserServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_ALREADY_EXISTS), "not created")
	}

	return &userapiv1.CreateUserServiceResponse{
		Success: true,
	}, status.Error(codes.Code(code.Code_OK), "created")
}

func (s *UserServiceServer) GetAllUsersService(ctx context.Context, req *userapiv1.GetAllUsersServiceRequest) (*userapiv1.GetAllUsersServiceResponse, error) {
	users, err := s.userSrv.All(ctx)
	if err != nil {
		return &userapiv1.GetAllUsersServiceResponse{},
			status.Error(codes.Code(code.Code_INTERNAL), "not read")
	}

	result := make([]*userapiv1.User, len(users))
	for i, u := range users {
		result[i] = &userapiv1.User{
			Id:             string(u.ID),
			Username:       string(u.Username),
			TimeOfCreation: u.TimeOfCreation.String(),
			Role:           string(u.Role),
			City:           u.City,
			Version:        int32(u.Version),
		}
	}
	return &userapiv1.GetAllUsersServiceResponse{
		Users: result,
	}, status.Error(codes.Code(code.Code_OK), "read")
}

func (s *UserServiceServer) GetInfoService(ctx context.Context, req *userapiv1.GetInfoServiceRequest) (*userapiv1.GetInfoServiceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &userapiv1.GetInfoServiceResponse{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	token, found := md["authorization"]
	if !found {
		return &userapiv1.GetInfoServiceResponse{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	username, err := s.authSrv.Username(model.JwtToken(token[0]))
	if err != nil {
		return &userapiv1.GetInfoServiceResponse{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	u, err := s.userSrv.ReadByUsername(ctx, username)
	if err != nil {
		return &userapiv1.GetInfoServiceResponse{}, status.Error(codes.Code(code.Code_INTERNAL), "not read")
	}

	return &userapiv1.GetInfoServiceResponse{
		User: &userapiv1.User{
			Id:             string(u.ID),
			Username:       string(u.Username),
			TimeOfCreation: u.TimeOfCreation.String(),
			Role:           string(u.Role),
			City:           u.City,
			Version:        int32(u.Version),
		},
	}, status.Error(codes.Code(code.Code_OK), "read")
}

func (s *UserServiceServer) UpdateByIdService(ctx context.Context, req *userapiv1.UpdateByIdServiceRequest) (*userapiv1.UpdateByIdServiceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &userapiv1.UpdateByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	token, found := md["authorization"]
	if !found {
		return &userapiv1.UpdateByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	role, err := s.authSrv.Role(model.JwtToken(token[0]))
	if err != nil {
		return &userapiv1.UpdateByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	if role != "admin" && role != "staff" {
		return &userapiv1.UpdateByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "not allowed")
	}

	err = s.userSrv.UpdateByID(ctx, model.ID(req.Id), model.UserInfo{
		Username: model.Username(req.Username),
		Password: model.Password(req.Password),
		Role:     model.Role(req.Role),
		City:     req.City,
		Version:  int32(req.Version),
	})
	if err != nil {
		return &userapiv1.UpdateByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_INTERNAL), "not updated")
	}

	return &userapiv1.UpdateByIdServiceResponse{
		Success: true,
	}, status.Error(codes.Code(code.Code_OK), "updated")

}

func (s *UserServiceServer) DeleteByIdService(ctx context.Context, req *userapiv1.DeleteByIdServiceRequest) (*userapiv1.DeleteByIdServiceResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &userapiv1.DeleteByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	token, found := md["authorization"]
	if !found {
		return &userapiv1.DeleteByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	role, err := s.authSrv.Role(model.JwtToken(token[0]))
	if err != nil {
		return &userapiv1.DeleteByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "missing token")
	}

	if role != "admin" {
		return &userapiv1.DeleteByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_PERMISSION_DENIED), "not allowed")
	}

	err = s.userSrv.DeleteByID(ctx, model.ID(req.Id))
	if err != nil {
		return &userapiv1.DeleteByIdServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_INTERNAL), "not deleted")
	}

	return &userapiv1.DeleteByIdServiceResponse{
		Success: true,
	}, status.Error(codes.Code(code.Code_OK), "deleted")

}
