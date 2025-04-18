package authusecase

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/logger"
)

type Signin interface {
	Run(ctx context.Context, payload *dtos.SigninRequest) (*dtos.SigninResponse, error)
}

type signinRepository interface {
	CheckUserPassword(ctx context.Context, user *entities.User) (*entities.User, error)
}

type signin struct {
	log  logger.Logger
	repo signinRepository
}

func NewSignin(ctx context.Context, log logger.Logger, repo signinRepository) Signin {
	return &signin{
		log:  log,
		repo: repo,
	}
}

func (s *signin) Run(ctx context.Context, payload *dtos.SigninRequest) (*dtos.SigninResponse, error) {
	//TODO implement me
	panic("implement me")
}
