package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
)

type SignupUseCase struct {
	log  logger.Logger
	repo ports.SignupRepository
}

func NewSignup(ctx context.Context, log logger.Logger, repo ports.SignupRepository) *SignupUseCase {
	return &SignupUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (*entities.SignupResponse, error) {
	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return nil, err
	}

	if exists.Status {
		var err = &apperrors.ErrLoginExists{}
		if exists.UsernameExists {
			err.Username = payload.Username
		}
		if exists.EmailExists {
			err.Email = payload.Email
		}

		return nil, err
	}

	return nil, nil
}
