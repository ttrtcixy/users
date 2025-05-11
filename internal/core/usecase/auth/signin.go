package authusecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/logger"
)

type SigninUseCase struct {
	log  logger.Logger
	repo usecaseports.SigninRepository
}

func NewSignin(ctx context.Context, log logger.Logger, repo usecaseports.Repository) *SigninUseCase {
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

func (u *SigninUseCase) comparePassword(hashedPassword, password []byte) bool {
	if len(hashedPassword) != len(password) {
		return false
	}

	for i, v := range hashedPassword {
		if v != password[i] {
			return false
		}
	}
	return true
}
