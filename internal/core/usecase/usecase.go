package usecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/usecase/auth"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/service/smtp"
)

type UseCase struct {
	*AuthUseCase
	*UserUseCase
}

func NewUseCase(ctx context.Context, log logger.Logger, repo usecaseports.Repository, cfg *config.Config, smtp smtp.Smtp) *UseCase {
	return &UseCase{
		NewAuthUseCase(ctx, log, repo, cfg, smtp),
		NewUserUseCase(ctx, log, repo),
	}
}

type AuthUseCase struct {
	*authusecase.SignoutUseCase
	*authusecase.SignupUseCase
	*authusecase.SigninUseCase
	*authusecase.VerifyUseCase
}

func NewAuthUseCase(ctx context.Context, log logger.Logger, repo usecaseports.Repository, cfg *config.Config, smtp smtp.Smtp) *AuthUseCase {
	return &AuthUseCase{
		SignoutUseCase: authusecase.NewSignout(ctx, log, repo),
		SignupUseCase:  authusecase.NewSignup(ctx, cfg.UsecaseConfig, log, repo, smtp),
		SigninUseCase:  authusecase.NewSignin(ctx, log, repo),
		VerifyUseCase:  authusecase.NewVerify(ctx, cfg.UsecaseConfig, log, repo),
	}
}

type UserUseCase struct {
}

func NewUserUseCase(ctx context.Context, log logger.Logger, repo usecaseports.Repository) *UserUseCase {
	return &UserUseCase{}
}
