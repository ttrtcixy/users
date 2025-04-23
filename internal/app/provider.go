package app

import (
	"context"
	"github.com/ttrtcixy/users/internal/app/closer"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/delivery/grpc"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/repository"
	"github.com/ttrtcixy/users/internal/storage"
	"github.com/ttrtcixy/users/internal/usecase"
	"time"
)

type Provider struct {
	logger logger.Logger
	cfg    *config.Config

	db storage.DB

	useCase *usecase.UseCase

	repository repository.Repository

	grpcServer *grpc.Server

	closer closer.Closer
}

func NewProvider() *Provider {
	return &Provider{}
}

func (p *Provider) UseCase() *usecase.UseCase {
	if p.useCase == nil {
		p.useCase = usecase.NewUseCase(context.Background(), p.Logger(), p.Repository())
	}

	return p.useCase
}

func (p *Provider) Logger() logger.Logger {
	if p.logger == nil {
		p.logger = logger.Load()
	}

	return p.logger
}

func (p *Provider) Config() *config.Config {
	if p.cfg == nil {
		p.cfg = config.NewConfig(p.Logger())

		p.Closer().Add(
			"env clear",
			p.cfg.Close,
		)
	}

	return p.cfg
}

func (p *Provider) Repository() repository.Repository {
	if p.repository == nil {
		p.repository = repository.NewRepository(context.Background(), p.Logger(), p.DB())
	}
	return p.repository
}

func (p *Provider) DB() storage.DB {
	if p.db == nil {
		p.db = storage.NewDB(context.Background(), p.Logger(), p.Config().DBConfig)

		p.Closer().Add(
			"db connect close",
			p.db.Close,
		)
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
