package usecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/repository"
	"github.com/ttrtcixy/users/internal/usecase/auth"
)

type UseCase struct {
	*AuthUseCase
}

func NewUseCase(ctx context.Context, log logger.Logger, repo repository.Repository) *UseCase {
	return &UseCase{
		NewAuthUseCase(ctx, log, repo),
	}
}

type AuthUseCase struct {
	authusecase.Signout
	authusecase.Signup
	authusecase.Signin
}

func NewAuthUseCase(ctx context.Context, log logger.Logger, repo repository.Repository) *AuthUseCase {
	return &AuthUseCase{
		Signout: authusecase.NewSignout(ctx, log, repo),
		Signup:  authusecase.NewSignup(ctx, log, repo),
		Signin:  authusecase.NewSignin(ctx, log, repo),
	}
}

type UserUseCase struct {
}
