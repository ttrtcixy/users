package usecaseports

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
	//User(ctx context.Context, user *entities.User) (*entities.User, error)
}

type VerifyRepository interface {
}
