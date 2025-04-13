package grpcauthservise

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type Signup interface {
	Signup(ctx context.Context, payload *dtos.SignupRequest) (*dtos.SignupResponse, error)
}

type signup struct {
	log     logger.Logger
	usecase *usecase.UseCase
}

func NewSignup(log logger.Logger, usecase *usecase.UseCase) Signup {
	return &signup{
		log:     log,
		usecase: usecase,
	}
}

func (s *signup) Signup(ctx context.Context, payload *dtos.SignupRequest) (*dtos.SignupResponse, error) {
	return nil, nil
}
