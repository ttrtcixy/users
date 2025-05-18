package ports

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
)

type Repository interface {
	AuthRepository
}

type AuthRepository interface {
	SignupRepository
	SigninRepository
	VerifyRepository
}

type SignupRepository interface {
	CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error)
	CreateUser(ctx context.Context, payload *entities.CreateUserRequest) error
	RunInTx(ctx context.Context, fn func(context.Context) error) error
}

type SigninRepository interface {
	UserByEmail(ctx context.Context, email string) (*entities.User, error)
	UserByUsername(ctx context.Context, username string) (*entities.User, error)
	CreateSession(ctx context.Context, payload *entities.CreateSession) error
}

type VerifyRepository interface {
	ActivateUser(ctx context.Context, email string) (*entities.User, error)
	CreateSession(ctx context.Context, payload *entities.CreateSession) error
}
