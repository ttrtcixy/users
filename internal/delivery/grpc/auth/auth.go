package grpcauthservise

import (
	"context"
	usersProtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/delivery/grpc/ports"
	"github.com/ttrtcixy/users/internal/logger"
)

type UserAuthService struct {
	*SigninService
	*SignupService
	*SignoutService
	usersProtos.UnsafeUsersAuthServer
}

func NewUserAuthService(ctx context.Context, log logger.Logger, usecase ports.UseCase) usersProtos.UsersAuthServer {
	return &UserAuthService{
		SigninService:  NewSignin(log, usecase),
		SignupService:  NewSignup(log, usecase),
		SignoutService: NewSignout(log, usecase),
	}
}
