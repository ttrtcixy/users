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

type VerifyUseCase interface {
	Verify(ctx context.Context, payload *entities.VerifyRequest) (*entities.VerifyResponse, error)
}

type RefreshUseCase interface {
	Refresh(ctx context.Context, payload *entities.RefreshRequest) (*entities.RefreshResponse, error)
}

type AuthUseCase interface {
	SigninUseCase
	SignupUseCase
	SignoutUseCase
	VerifyUseCase
	RefreshUseCase
}

type UseCase interface {
	AuthUseCase
}
