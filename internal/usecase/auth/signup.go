package authusecase

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
)

type Signup interface {
	Run(ctx context.Context, payload *dtos.SignupRequest) (*dtos.SignupResponse, error)
}

type signupRepository interface {
}

type signup struct {
	log  logger.Logger
	repo signupRepository
}

func NewSignup(ctx context.Context, log logger.Logger, repo signupRepository) Signup {
	return &signup{
		log:  log,
		repo: repo,
	}
}

func (l *signup) Run(ctx context.Context, payload *dtos.SignupRequest) (*dtos.SignupResponse, error) {
	//TODO implement me
	panic("implement me")
}
