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

type VerifyUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo ports.VerifyRepository
	jwt  ports.JwtService
	hash ports.HasherService
}

type VerifyUseCaseDependency struct {
	Cfg  *config.UsecaseConfig
	Log  logger.Logger
	Repo ports.VerifyRepository
	Jwt  ports.JwtService
	Hash ports.HasherService
}

func NewVerify(ctx context.Context, dep *VerifyUseCaseDependency) *VerifyUseCase {
	return &VerifyUseCase{
		cfg:  dep.Cfg,
		log:  dep.Log,
		repo: dep.Repo,
		jwt:  dep.Jwt,
		hash: dep.Hash,
	}
}

// Verify - get jwtToken with email and activate user with that email.
func (u *VerifyUseCase) Verify(ctx context.Context, payload *entities.VerifyRequest) (result *entities.VerifyResponse, err error) {
	const op = "VerifyUseCase.Verify"
	defer func() {
		if err != nil {
			if errors.Is(err, apperrors.ErrEmailTokenExpired) {
				return
			}
			u.log.ErrorOp(op, err)
			err = apperrors.ErrServer
		}
	}()

	user, err := u.activateUser(ctx, payload.JwtToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwt.AccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, err := u.createSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result = &entities.VerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}

func (u *VerifyUseCase) activateUser(ctx context.Context, jwtToken string) (user *entities.User, err error) {
	const op = "activateUser"

	email, err := u.jwt.ParseVerificationToken(jwtToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if user, err = u.repo.ActivateUser(ctx, email); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *VerifyUseCase) createSession(ctx context.Context, userID int64) (refreshToken string, err error) {
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

	return refreshTokenHash, nil
}
