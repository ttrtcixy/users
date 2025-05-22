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

type RefreshUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo ports.RefreshRepository
	jwt  *token.JwtTokenService
	//hash ports.HasherService
}

type RefreshUseCaseDependency struct {
	Cfg  *config.UsecaseConfig
	Log  logger.Logger
	Repo ports.RefreshRepository
	Jwt  *token.JwtTokenService
	//Hash ports.HasherService
}

func NewRefresh(ctx context.Context, dep *RefreshUseCaseDependency) *RefreshUseCase {
	return &RefreshUseCase{
		cfg:  dep.Cfg,
		log:  dep.Log,
		repo: dep.Repo,
		jwt:  dep.Jwt,
		//hash: dep.Hash,
	}
}

func (u *RefreshUseCase) Refresh(ctx context.Context, payload *entities.RefreshRequest) (result *entities.RefreshResponse, err error) {
	const op = "RefreshUseCase.Refresh"
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

	// parse client token
	clientID, JTI, err := u.jwt.ParseRefreshToken(payload.RefreshToken)
	if err != nil {
		return nil, err
	}

	// new refresh token info refresh token
	exp := time.Now().Add(u.cfg.RefreshJwtExpiry())
	newJti := uuid.NewString()

	createReq := &entities.RefreshSession{
		ClientUUID:          clientID,
		OldRefreshTokenUUID: JTI,
		NewRefreshTokenUUID: newJti,
		ExpiresAt:           exp,
	}

	// check token jtl + check clientID if good add new refresh jti and return public user info
	userInfo, err := u.repo.RefreshSession(ctx, createReq)
	if err != nil {
		return nil, err
	}

	// new refresh token
	refreshToken, err := u.jwt.RefreshToken(clientID, newJti, exp)
	if err != nil {
		return nil, err
	}
	// new access token
	accessToken, err := u.jwt.AccessToken(userInfo)

	return &entities.RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientID:     clientID,
	}, nil
}
