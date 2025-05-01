package authusecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/entities"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/ports"
)

type SignupUseCase struct {
	log  logger.Logger
	repo ports.SignupRepository
	cfg  config.UsecaseConfig
}

func NewSignup(ctx context.Context, log logger.Logger, repo ports.SignupRepository, cfg config.UsecaseConfig) *SignupUseCase {
	return &SignupUseCase{
		log:  log,
		repo: repo,
		cfg:  cfg,
	}
}

func (u *SignupUseCase) Signup(ctx context.Context, payload *entities.SignupRequest) (*entities.SignupResponse, error) {
	exists, err := u.repo.CheckLoginExist(ctx, payload)
	if err != nil {
		return nil, err
	}

	if exists.Status {
		var err = &apperrors.ErrLoginExists{}
		if exists.UsernameExists {
			err.Username = payload.Username
		}
		if exists.EmailExists {
			err.Email = payload.Email
		}
		return nil, err
	}

	hash, err := u.hash(payload.Password)
	if err != nil {
		return nil, err
	}
	_ = hash
	return nil, nil
}

// hash generates a hash of the password using a salt
func (u *SignupUseCase) hash(password string) (string, error) {
	salt, err := u.salt()
	if err != nil {
		return "", err
	}

	hasher := sha256.New()

	bytePassword := append([]byte(password), salt...)
	hasher.Write(bytePassword)

	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash), nil
}

// salt generate random salt for password
func (u *SignupUseCase) salt() ([]byte, error) {
	salt := make([]byte, u.cfg.PasswordSaltLength())
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}
