package grpcauthservise

import (
	"context"
	usersProtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type UserAuthService struct {
	Signin
	Signup
	Signout
	usersProtos.UnsafeUsersAuthServer
}

func NewUserAuthService(ctx context.Context, log logger.Logger, usecase *usecase.UseCase) usersProtos.UsersAuthServer {
	return &UserAuthService{
		Signin:  NewSignin(log, usecase),
		Signup:  NewSignup(log, usecase),
		Signout: NewSignout(log, usecase),
	}
}

//func (s *UserAuthService) Signup(ctx context.Context, payload *dtos.SignupRequest) (*dtos.SignupResponse, error) {
//	return nil, nil
//}

//func (s *UserAuthService) Signin(ctx context.Context, payload *dtos.SigninRequest) (*dtos.SigninResponse, error) {
//	return nil, nil
//}

//func (s *UserAuthService) Signout(ctx context.Context, payload *dtos.SignoutRequest) (*emptypb.Empty, error) {
//	return nil, nil
//}
