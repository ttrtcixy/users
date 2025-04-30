package authusecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
	"strings"
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
	// todo check username or email exists
	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return nil, err
	}

	if exists.Status {
		var str strings.Builder
		if exists.UsernameExists {
			str.WriteString(fmt.Sprintf("username: %s, exists; ", payload.Username))
		}
		if exists.EmailExists {
			str.WriteString(fmt.Sprintf("email: %s, exists; ", payload.Email))
		}

		return nil, errors.New(str.String())
	}

	// todo hash password
	// todo add new user
	// todo check user email
	// todo generate and send tokens

	return nil, nil
}
