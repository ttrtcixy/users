package grpcauthservise

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type AuthSignin interface {
	Signin(ctx context.Context, payload *dtos.SigninRequest) (*dtos.SigninResponse, error)
}

type signin struct {
	log     logger.Logger
	usecase *usecase.UseCase
}

func NewSignin(log logger.Logger, usecase *usecase.UseCase) AuthSignin {
	return &signin{
		log:     log,
		usecase: usecase,
	}
}

func (s *signin) Signin(ctx context.Context, payload *dtos.SigninRequest) (*dtos.SigninResponse, error) {
	return nil, nil
}
