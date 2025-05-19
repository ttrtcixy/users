package provider

import (
	"context"
	"github.com/ttrtcixy/users/internal/app/closer"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/repository"
	"github.com/ttrtcixy/users/internal/core/usecase"
	"github.com/ttrtcixy/users/internal/delivery/grpc"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/service/hash"
	token "github.com/ttrtcixy/users/internal/service/jwt"
	"github.com/ttrtcixy/users/internal/service/smtp"
	storage "github.com/ttrtcixy/users/internal/storage/pg"
)

type Provider struct {
	logger logger.Logger
	closer closer.Closer

	cfg *config.Config

	db   storage.DB
	smtp *smtp.SenderService
	hash *hash.HasherService
	jwt  *token.JwtTokenService

	useCase *usecase.UseCase

	repository *repository.Repository

	grpcServer *grpc.Server
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) UseCase() *usecase.UseCase {
	if p.useCase == nil {
		dep := &usecase.Dependency{
			Cfg:  p.Config(),
			Log:  p.Logger(),
			Repo: p.Repository(),
			Smtp: p.SmtpService(),
			Hash: p.HasherService(),
			Jwt:  p.JwtService(),
		}
		p.useCase = usecase.NewUseCase(context.Background(), dep)
	}

	return p.useCase
}

func (p *Provider) Logger() logger.Logger {
	if p.logger == nil {
		p.logger = logger.Load()
	}

	return p.logger
}

func (p *Provider) SmtpService() *smtp.SenderService {
	if p.smtp == nil {
		p.smtp = smtp.New(p.Config().SmtpConfig)
	}

	return p.smtp
}

func (p *Provider) Config() *config.Config {
	if p.cfg == nil {
		cfg, err := config.New()
		if err != nil {
			p.Logger().Fatal("[!] %s", err.Error())
		}

		p.cfg = cfg
		p.Closer().Add(
			"env clear",
			p.cfg.Close,
		)

		p.Logger().Info("[+] config loaded")
	}

	return p.cfg
}

func (p *Provider) Repository() *repository.Repository {
	if p.repository == nil {
		p.repository = repository.NewRepository(context.Background(), p.Logger(), p.DB())
	}
	return p.repository
}

func (p *Provider) DB() storage.DB {
	if p.db == nil {
		db, err := storage.New(context.Background(), p.Logger(), p.Config().DBConfig)
		if err != nil {
			p.Logger().Error("[!] %s", err.Error())
			p.Closer().Close()
		}

		p.db = db
		p.Closer().Add(
			"db connect close",
			p.db.Close,
		)
		p.Logger().Info("[+] db connected")
	}
	return p.db
}

func (p *Provider) GRPCServer() *grpc.Server {
	if p.grpcServer == nil {
		p.grpcServer = grpc.NewGRPCServer(p.Logger(), p.Config().GRPCServerConfig, p.UseCase())
		p.Closer().Add(
			"grpc server close",
			p.grpcServer.Close,
		)
		p.Logger().Info("[+] grpc server started")
	}
	return p.grpcServer
}

func (p *Provider) Closer() closer.Closer {
	if p.closer == nil {
		p.closer = closer.New(closer.Config{
			TotalDuration: p.Config().CloserConfig.TotalDuration(),
			FuncDuration:  p.Config().CloserConfig.FuncDuration(),
			Logger:        p.Logger(),
		})
	}
	return p.closer
}

func (p *Provider) JwtService() *token.JwtTokenService {
	if p.jwt == nil {
		p.jwt = token.New(p.Config().UsecaseConfig.JWTConfig)
	}

	return p.jwt
}

func (p *Provider) HasherService() *hash.HasherService {
	if p.hash == nil {
		p.hash = hash.New()
	}
	return p.hash
}
