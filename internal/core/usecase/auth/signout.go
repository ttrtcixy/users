package authusecase

import (
	"context"
	"errors"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
)

type SignoutUseCase struct {
	log  logger.Logger
	repo ports.SignoutRepository
}

func NewSignout(ctx context.Context, log logger.Logger, repo ports.Repository) *SignoutUseCase {
	return &SignoutUseCase{
		log:  log,
		repo: repo,
	}
}

func (u *SignoutUseCase) Signout(ctx context.Context, payload *entities.SignoutRequest) (err error) {
	const op = "SignoutUseCase.Signout"

	defer func() {
		if err != nil {
			var userErr apperrors.UserError
			if errors.As(err, &userErr) {
				return
			}

			u.log.ErrorOp(op, err)
			err = apperrors.ErrServer
		}
	}()

	if err := u.repo.DeleteSession(ctx, payload); err != nil {
		return apperrors.Wrap(op, err)
	}

	return nil
}
