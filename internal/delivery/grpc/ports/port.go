package ports

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
)

type SigninUseCase interface {
	Signin(ctx context.Context, payload *entities.SigninRequest) (*entities.SigninResponse, error)
}

type SignupUseCase interface {
	Signup(ctx context.Context, payload *entities.SignupRequest) error
}

type SignoutUseCase interface {
	Signout(ctx context.Context, payload *entities.SignoutRequest) error
}

type AuthUseCase interface {
	SigninUseCase
	SignupUseCase
	SignoutUseCase
}

type UseCase interface {
	AuthUseCase
}
