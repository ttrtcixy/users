package authusecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
)

type SignupUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo ports.SignupRepository
	smtp ports.SmtpService
	hash ports.HasherService
	jwt  ports.JwtService
}

type SignupUseCaseDeps struct {
	Cfg  *config.UsecaseConfig
	Log  logger.Logger
	Repo ports.SignupRepository
	Smtp ports.SmtpService
	Hash ports.HasherService
	Jwt  ports.JwtService
}

func NewSignup(ctx context.Context, dep *SignupUseCaseDeps) *SignupUseCase {
	return &SignupUseCase{
		cfg:  dep.Cfg,
		log:  dep.Log,
		repo: dep.Repo,
		smtp: dep.Smtp,
		hash: dep.Hash,
		jwt:  dep.Jwt,
	}
}

// todo validate как работает fmt.Errorf("%s: %w", на большом стеке вызова)

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (err error) {
	const op = "SignupUseCase.Signup"

	err = u.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.validPayload(ctx, payload); err != nil {
			return err
		}

		hash, salt, err := u.passwordHashing(payload.Password)
		if err != nil {
			return err
		}

		token, err := u.jwt.EmailVerificationToken(payload.Email)
		if err != nil {
			return err
		}

		createReq := &entities.CreateUserRequest{
			Username:     payload.Username,
			Email:        payload.Email,
			PasswordHash: hash,
			PasswordSalt: salt,
		}
		if err = u.repo.CreateUser(ctx, createReq); err != nil {
			return err
		}

		if err = u.smtp.DebugSend(payload.Email, token); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		var ue apperrors.UserError
		if errors.As(err, &ue) {
			return err
		}

		u.log.ErrorOp(op, err)
		return apperrors.ErrServer
	}
	return nil
}

// validPayload - UserError: *apperrors.ErrLoginExists
func (u *SignupUseCase) validPayload(ctx context.Context, payload *entities.SignupRequest) error {
	const op = "validPayload"

	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if exists.Status {
		var err = &apperrors.ErrLoginExists{}
		if exists.UsernameExists {
			err.Username = payload.Username
		}
		if exists.EmailExists {
			err.Email = payload.Email
		}
		return err
	}

	return nil
}

func (u *SignupUseCase) passwordHashing(password string) (hash string, salt string, err error) {
	const op = "passwordHashing"

	byteSalt, err := u.hash.Salt(u.cfg.PasswordSaltLength())
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if hash, err = u.hash.HashWithSalt(password, byteSalt); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return hash, base64.StdEncoding.EncodeToString(byteSalt), nil
}
