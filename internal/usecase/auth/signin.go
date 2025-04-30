package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
)

type SigninUseCase struct {
	log  logger.Logger
	repo ports.Repository
}

func NewSignin(ctx context.Context, log logger.Logger, repo ports.Repository) *SigninUseCase {
	return &SigninUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *SigninUseCase) Signin(ctx context.Context, payload *entities.SigninRequest) (*entities.SigninResponse, error) {
	// todo check signin method(email/username)
	// todo get user id by method + password
	// todo generate and send tokens
	return nil, nil
}
