package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/smtp"
	"time"
)

type SignupUseCase struct {
	cfg  *config.UsecaseConfig
	log  logger.Logger
	repo usecaseports.SignupRepository
	smtp smtp.Smtp
}

func NewSignup(ctx context.Context, cfg *config.UsecaseConfig, log logger.Logger, repo usecaseports.SignupRepository, smtp smtp.Smtp) *SignupUseCase {
	return &SignupUseCase{
		cfg:  cfg,
		log:  log,
		repo: repo,
		smtp: smtp,
	}
}

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (err error) {
	return u.repo.RunInTx(ctx, func(ctx context.Context) error {
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
			Username:        payload.Username,
			Email:           payload.Email,
			PasswordHash:    hash,
			PasswordSalt:    salt,
			ActivationToken: token,
		}
		if err = u.repo.CreateUser(ctx, createReq); err != nil {
			return err
		}

		if err = u.smtp.DebugSend(payload.Email, token); err != nil {
			return err
		}

		return nil
	})
}

func (u *SignupUseCase) validPayload(ctx context.Context, payload *entities.SignupRequest) error {
	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return err
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
	byteSalt, err := u.salt(u.cfg.PasswordSaltLength())
	if err != nil {
		return "", "", err
	}

	if hash, err = u.hashWithSalt(password, byteSalt); err != nil {
		return "", "", err
	}

	return hash, base64.StdEncoding.EncodeToString(byteSalt), nil
}

// salt generate random salt
func (u *SignupUseCase) salt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// hash generates a hash with salt using sha256 and return base64 string
func (u *SignupUseCase) hashWithSalt(str string, salt []byte) (hash string, err error) {
	hasher := sha256.New()
	data := append([]byte(str), salt...)
	if _, err = hasher.Write(data); err != nil {
		return "", fmt.Errorf("hashWithSalt: write failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil

}

func (u *SignupUseCase) jwt(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"email": email, "exp": time.Now().Add(time.Hour * 24)})

	return token.SignedString(u.cfg.JWTSecret())
}

// hash generates a hash using sha256 and return base64 string
//func (u *SignupUseCase) hash(str string) (hash string, err error) {
//	hasher := sha256.New()
//	if _, err = hasher.Write([]byte(str)); err != nil {
//		return "", fmt.Errorf("hash: write failed: %w", err)
//	}
//	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
//}

//func (u *SignupUseCase) token() (token string, err error) {
//	var rnd = make([]byte, 32)
//	if _, err = rand.Read(rnd); err != nil {
//		return "", err
//	}
//	return base64.URLEncoding.EncodeToString(rnd), nil
//}
