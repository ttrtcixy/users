package usecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase/auth"
	"github.com/ttrtcixy/users/internal/usecase/ports"
)

type UseCase struct {
	*AuthUseCase
	*UserUseCase
}

func NewUseCase(ctx context.Context, log logger.Logger, repo ports.Repository, cfg *config.Config) *UseCase {
	return &UseCase{
		NewAuthUseCase(ctx, log, repo, cfg),
		NewUserUseCase(ctx, log, repo),
	}
}

type AuthUseCase struct {
	*authusecase.SignoutUseCase
	*authusecase.SignupUseCase
	*authusecase.SigninUseCase
}

func NewAuthUseCase(ctx context.Context, log logger.Logger, repo ports.Repository, cfg *config.Config) *AuthUseCase {
	return &AuthUseCase{
		SignoutUseCase: authusecase.NewSignout(ctx, log, repo),
		SignupUseCase:  authusecase.NewSignup(ctx, log, cfg.UsecaseConfig, repo),
		SigninUseCase:  authusecase.NewSignin(ctx, log, repo),
	}
}

type UserUseCase struct {
}

func NewUserUseCase(ctx context.Context, log logger.Logger, repo ports.Repository) *UserUseCase {
	return &UserUseCase{}
}
