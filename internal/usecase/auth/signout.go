package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
)

type SignoutUseCase struct {
	log  logger.Logger
	repo ports.Repository
}

func NewSignout(ctx context.Context, log logger.Logger, repo ports.Repository) *SignoutUseCase {
	return &SignoutUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *SignoutUseCase) Signout(ctx context.Context, payload *entities.SignoutRequest) error {
	return nil
}
