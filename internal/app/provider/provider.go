package provider

import (
	"context"
	"github.com/ttrtcixy/users/internal/app/closer"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/core/repository"
	"github.com/ttrtcixy/users/internal/core/usecase"
	"github.com/ttrtcixy/users/internal/delivery/grpc"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/smtp"
	storage "github.com/ttrtcixy/users/internal/storage/pg"
	"time"
)

type Provider struct {
	logger logger.Logger
	closer closer.Closer

	cfg *config.Config

	db   storage.DB
	smtp smtp.Smtp

	useCase *usecase.UseCase

	repository *repository.Repository

	grpcServer *grpc.Server
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) UseCase() *usecase.UseCase {
	if p.useCase == nil {
		p.useCase = usecase.NewUseCase(context.Background(), p.Logger(), p.Repository(), p.Config(), p.Smtp())
	}

	return p.useCase
}

func (p *Provider) Logger() logger.Logger {
	if p.logger == nil {
		p.logger = logger.Load()
	}

	return p.logger
}

func (p *Provider) Smtp() smtp.Smtp {
	if p.smtp == nil {
		p.smtp = smtp.New(p.Config().SmtpConfig)
	}

	return p.smtp
}

func (p *Provider) Config() *config.Config {
	if p.cfg == nil {
		cfg, err := config.New()
		if err != nil {
			p.Logger().Error("[!] %s", err.Error())
			p.Closer().Close()
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
			TotalDuration: 5 * time.Second,
			FuncDuration:  3 * time.Second,
			Logger:        p.Logger(),
		})
	}
	return p.closer
}
