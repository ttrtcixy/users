package usecase

import (
	"context"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/usecase/auth"
	"github.com/ttrtcixy/users/internal/core/usecase/ports"
	"github.com/ttrtcixy/users/internal/logger"
)

type UseCase struct {
	*AuthUseCase
	*UserUseCase
}

func NewUseCase(ctx context.Context, dep *Dependency) *UseCase {
	return &UseCase{
		NewAuthUseCase(ctx, dep),
		NewUserUseCase(ctx, dep),
	}
}

type AuthUseCase struct {
	*authusecase.SignoutUseCase
	*authusecase.SignupUseCase
	*authusecase.SigninUseCase
	*authusecase.VerifyUseCase
}

type Dependency struct {
	Cfg  *config.Config
	Log  logger.Logger
	Repo ports.Repository
	Smtp ports.SmtpService
	Hash ports.HasherService
	Jwt  ports.JwtService
}

func NewAuthUseCase(ctx context.Context, dep *Dependency) *AuthUseCase {
	return &AuthUseCase{
		SignoutUseCase: authusecase.NewSignout(ctx, dep.Log, dep.Repo),
		SignupUseCase: authusecase.NewSignup(ctx, &authusecase.SignupUseCaseDeps{
			Cfg:  dep.Cfg.UsecaseConfig,
			Log:  dep.Log,
			Repo: dep.Repo,
			Smtp: dep.Smtp,
			Hash: dep.Hash,
			Jwt:  dep.Jwt,
		}),
		SigninUseCase: authusecase.NewSignin(ctx, &authusecase.SigninUseCaseDeps{
			Cfg:  dep.Cfg.UsecaseConfig,
			Log:  dep.Log,
			Repo: dep.Repo,
			Jwt:  dep.Jwt,
			Hash: dep.Hash,
		}),
		VerifyUseCase: authusecase.NewVerify(ctx, &authusecase.VerifyUseCaseDependency{
			Cfg:  dep.Cfg.UsecaseConfig,
			Log:  dep.Log,
			Repo: dep.Repo,
			Jwt:  dep.Jwt,
			Hash: dep.Hash,
		}),
	}
}

type UserUseCase struct {
}

func NewUserUseCase(ctx context.Context, dep *Dependency) *UserUseCase {
	return &UserUseCase{}
}
