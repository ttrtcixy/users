package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/logger"
)

type SignoutUseCase struct {
	log  logger.Logger
	repo usecaseports.Repository
}

func NewSignout(ctx context.Context, log logger.Logger, repo usecaseports.Repository) *SignoutUseCase {
	return &SignoutUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *SignoutUseCase) Signout(ctx context.Context, payload *entities.SignoutRequest) error {
	return nil
}
