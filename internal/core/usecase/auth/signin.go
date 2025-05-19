package authusecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"time"
)

type SigninUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo ports.SigninRepository
	jwt  ports.JwtService
	hash ports.HasherService
}

type SigninUseCaseDeps struct {
	Cfg  *config.UsecaseConfig
	Log  logger.Logger
	Repo ports.SigninRepository
	Jwt  ports.JwtService
	Hash ports.HasherService
}

func NewSignin(ctx context.Context, dep *SigninUseCaseDeps) *SigninUseCase {
	return &SigninUseCase{
		cfg:  dep.Cfg,
		log:  dep.Log,
		repo: dep.Repo,
		jwt:  dep.Jwt,
		hash: dep.Hash,
	}
}

// todo max session count
func (u *SigninUseCase) Signin(ctx context.Context, payload *entities.SigninRequest) (result *entities.SigninResponse, err error) {
	const op = "SigninUseCase.Signin"
	defer func() {
		if err != nil {
			// todo userError interface
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

	user, err := u.validateUser(ctx, payload)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwt.AccessToken(user)
	if err != nil {
		return nil, err
	}

	// todo if refresh token exists, return him
	refreshToken, err := u.createSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result = &entities.SigninResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

func (u *SigninUseCase) validateUser(ctx context.Context, payload *entities.SigninRequest) (user *entities.User, err error) {
	const op = "validateUser"
	// todo test с использование не существующего пользователя
	if user, err = u.user(ctx, payload); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user.IsActive == false {
		return nil, apperrors.ErrEmailVerify
	}

	// todo test
	ok, err := u.hash.ComparePasswords(user.PasswordHash, payload.Password, user.PasswordSalt)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		return nil, apperrors.ErrInvalidPassword
	}

	return user, nil
}

func (u *SigninUseCase) user(ctx context.Context, payload *entities.SigninRequest) (user *entities.User, err error) {
	const op = "user"

	if payload.Email != "" {
		user, err = u.repo.UserByEmail(ctx, payload.Email)
	} else {
		user, err = u.repo.UserByUsername(ctx, payload.Username)
	}
	if err != nil {
		// todo refactor
		if errors.Is(err, apperrors.ErrUserNotRegister) {
			return nil, err
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// todo other service??
func (u *SigninUseCase) createSession(ctx context.Context, userID int64) (refreshToken string, err error) {
	const op = "createSession"

	clientUUID := uuid.NewString()

	tokenUUID := uuid.NewString()

	exp := time.Now().Add(u.cfg.RefreshJwtExpiry())

	if refreshToken, err = u.jwt.RefreshToken(clientUUID, tokenUUID, exp); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	refreshTokenHash, err := u.hash.Hash(refreshToken)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	createReq := &entities.CreateSession{
		UserID:           userID,
		RefreshTokenHash: refreshTokenHash,
		ClientUUID:       clientUUID,
		RefreshTokenUUID: tokenUUID,
		ExpiresAt:        exp,
	}

	if err = u.repo.CreateSession(ctx, createReq); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}
