package controller

import (
	"context"
	"user-management/internal/model"
	authapiv1 "user-management/internal/proto/usermgt/authapi/v1"
	service "user-management/internal/service"

	"google.golang.org/genproto/googleapis/rpc/code"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceServer struct {
	userSrv service.UserService
	authSrv service.AuthService
	authapiv1.UnimplementedAuthServiceServer
}

func NewAuthServiceServer(us service.UserService, as service.AuthService) *AuthServiceServer {
	return &AuthServiceServer{
		userSrv: us,
		authSrv: as,
	}
}

// var _ authapiv1.AuthServiceServer = (*AuthServiceServer)(nil)

func (s *AuthServiceServer) LoginService(ctx context.Context, req *authapiv1.LoginServiceRequest) (*authapiv1.LoginServiceResponse, error) {
	token, err := s.authSrv.Login(ctx, model.Username(req.Username), model.Password(req.Password))
	if err != nil {
		return &authapiv1.LoginServiceResponse{}, status.Error(codes.Code(code.Code_UNAUTHENTICATED), "login failed")
	}

	return &authapiv1.LoginServiceResponse{
		JwtToken: string(token),
	}, status.Error(codes.Code(code.Code_OK), "loged in")
}

func (s *AuthServiceServer) UpdatePasswordService(ctx context.Context, req *authapiv1.UpdatePasswordServiceRequest) (*authapiv1.UpdatePasswordServiceResponse, error) {
	err := s.authSrv.UpdatePassword(ctx, model.UpdatePassword{
		Username:    model.Username(req.Username),
		OldPassword: model.Password(req.OldPass),
		NewPassword: model.Password(req.NewPass),
	})
	if err != nil {
		return &authapiv1.UpdatePasswordServiceResponse{
			Success: false,
		}, status.Error(codes.Code(code.Code_UNKNOWN), "password not updated")
	}

	return &authapiv1.UpdatePasswordServiceResponse{
		Success: true,
	}, status.Error(codes.Code(code.Code_OK), "password updated")
}
