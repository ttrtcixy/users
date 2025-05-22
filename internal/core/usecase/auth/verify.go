package authusecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	token "github.com/ttrtcixy/users/internal/service/jwt"
	"time"
)

type VerifyUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo ports.VerifyRepository
	jwt  *token.JwtTokenService
	//hash ports.HasherService
}

type VerifyUseCaseDependency struct {
	Cfg  *config.UsecaseConfig
	Log  logger.Logger
	Repo ports.VerifyRepository
	Jwt  *token.JwtTokenService
	//Hash ports.HasherService
}

func NewVerify(ctx context.Context, dep *VerifyUseCaseDependency) *VerifyUseCase {
	return &VerifyUseCase{
		cfg:  dep.Cfg,
		log:  dep.Log,
		repo: dep.Repo,
		jwt:  dep.Jwt,
		//hash: dep.Hash,
	}
}

// Verify - get jwtToken with email and activate user with that email.
func (u *VerifyUseCase) Verify(ctx context.Context, payload *entities.VerifyRequest) (result *entities.VerifyResponse, err error) {
	const op = "VerifyUseCase.Verify"

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

	user, err := u.activateUser(ctx, payload.JwtEmailToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := u.jwt.AccessToken(user)
	if err != nil {
		return nil, err
	}

	refreshToken, clientUUID, err := u.createSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	result = &entities.VerifyResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientUUID:   clientUUID,
	}

	return result, nil
}

func (u *VerifyUseCase) activateUser(ctx context.Context, jwtToken string) (user *entities.TokenUserInfo, err error) {
	const op = "activateUser"

	email, err := u.jwt.ParseVerificationToken(jwtToken)
	if err != nil {
		return nil, apperrors.Wrap(op, err)
	}

	if user, err = u.repo.ActivateUser(ctx, email); err != nil {
		return nil, apperrors.Wrap(op, err)
	}

	return user, nil
}

func (u *VerifyUseCase) createSession(ctx context.Context, userID int) (refreshToken, clientID string, err error) {
	const op = "createSession"

	clientUUID := uuid.NewString()

	tokenUUID := uuid.NewString()

	exp := time.Now().Add(u.cfg.RefreshJwtExpiry())

	if refreshToken, err = u.jwt.RefreshToken(clientUUID, tokenUUID, exp); err != nil {
		return "", "", apperrors.Wrap(op, err)
	}

	createReq := &entities.CreateSession{
		UserID:           userID,
		ClientUUID:       clientUUID,
		RefreshTokenUUID: tokenUUID,
		ExpiresAt:        exp,
	}

	if err = u.repo.CreateSession(ctx, createReq); err != nil {
		return "", "", apperrors.Wrap(op, err)
	}

	return refreshToken, clientID, nil
}
