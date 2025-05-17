package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/service/smtp"
	"time"
)

type SignupUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo usecaseports.SignupRepository
	smtp smtp.Smtp
}

type EmailVerificationClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewSignup(ctx context.Context, cfg *config.UsecaseConfig, log logger.Logger, repo usecaseports.SignupRepository, smtp smtp.Smtp) *SignupUseCase {
	return &SignupUseCase{
		cfg:  cfg,
		log:  log,
		repo: repo,
		smtp: smtp,
	}
}

const op = "SignupUseCase.Signup"

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (err error) {
	err = u.repo.RunInTx(ctx, func(ctx context.Context) error {
		if err := u.validPayload(ctx, payload); err != nil {
			return err
		}

		hash, salt, err := u.passwordHashing(payload.Password)
		if err != nil {
			return err
		}

		token, err := u.jwt(payload.Email)
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

	byteSalt, err := u.salt(u.cfg.PasswordSaltLength())
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if hash, err = HashWithSalt(password, byteSalt); err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return hash, base64.StdEncoding.EncodeToString(byteSalt), nil
}

// salt generate random salt
func (u *SignupUseCase) salt(length int) ([]byte, error) {
	const op = "salt"

	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return salt, nil
}

// HashWithSalt generates a hash with salt using sha256 and return base64 string
func HashWithSalt(str string, salt []byte) (hash string, err error) {
	const op = "hashWithSalt"

	hasher := sha256.New()
	data := append([]byte(str), salt...)
	if _, err = hasher.Write(data); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func (u *SignupUseCase) jwt(email string) (token string, err error) {
	const op = "jwt"

	expAt := time.Now().Add(u.cfg.EmailJwtExpiry())

	claims := EmailVerificationClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "auth_grpc_app",
			ExpiresAt: jwt.NewNumericDate(expAt),
		},
	}

	if token, err = JWT(u.cfg.JWTSecret(), claims); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
