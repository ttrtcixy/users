package grpcauthservise

import (
	"context"
	usersProtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type UserAuthService struct {
	AuthSignin
	AuthSignup
	AuthSignout
	usersProtos.UnsafeUsersAuthServer
}

func NewUserAuthService(ctx context.Context, log logger.Logger, usecase *usecase.UseCase) usersProtos.UsersAuthServer {
	return &UserAuthService{
		AuthSignin:  NewSignin(log, usecase),
		AuthSignup:  NewSignup(log, usecase),
		AuthSignout: NewSignout(log, usecase),
	}
}
