package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/logger"
)

type VerifyUseCase struct {
	log  logger.Logger
	repo usecaseports.usecaseports
}

func NewVerify(ctx context.Context, log logger.Logger, repo usecaseports.Repository) *VerifyUseCase {
	return &VerifyUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *VerifyUseCase) Verify(ctx context.Context, payload *entities.VerifyRequest) (*entities.VerifyResponse, error) {
	// todo add tx
	// todo generate token
	token, err := u.token()
	if err != nil {
		return nil, err
	}
	_ = token
	// todo hash token

	// todo add hash in database
	// todo send token to user
	return nil, nil
}
