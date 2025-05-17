package authusecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
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

func (u *SigninUseCase) Signin(ctx context.Context, payload *entities.SigninRequest) (result *entities.SigninResponse, err error) {
	const op = "SigninUseCase.Signin"
	defer func() {
		if err != nil {
			if errors.Is(err, apperrors.ErrEmailVerify) {
				return
			}
			if errors.Is(err, apperrors.ErrUserNotRegister) {
				return
			}
			if errors.Is(err, apperrors.ErrInvalidPassword) {
				return
			}
			u.log.ErrorOp(op, err)
			err = apperrors.ErrServer
		}
	}()

	// todo test с использование не существующего пользователя
	user, err := u.user(ctx, payload)
	if err != nil {
		return nil, err
	}

	if user.IsActive == false {
		return nil, apperrors.ErrEmailVerify
	}

	if err = u.checkPassword(user.PasswordHash, payload.Password, user.PasswordSalt); err != nil {
		return nil, err
	}

	// todo generate and send tokens
	return nil, nil
}

func (u *SigninUseCase) user(ctx context.Context, payload *entities.SigninRequest) (user *entities.User, err error) {
	const op = "user"

	if payload.Email != "" {
		user, err = u.repo.UserByEmail(ctx, payload.Email)
	} else {
		user, err = u.repo.UserByUsername(ctx, payload.Username)
	}
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *SigninUseCase) checkPassword(hashedPassword, password, salt string) error {
	const op = "checkPassword"

	hash, err := HashWithSalt(password, []byte(salt))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if ok := u.comparePassword(hashedPassword, hash); !ok {
		return apperrors.ErrInvalidPassword
	}
	return nil
}

func (u *SigninUseCase) comparePassword(hashedPassword, password string) bool {
	return hashedPassword == password
}
