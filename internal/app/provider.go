package app

import (
	"context"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/repository"
	"github.com/ttrtcixy/users/internal/storage"
	"github.com/ttrtcixy/users/internal/usecase"
)

type Provider struct {
	logger logger.Logger
	cfg    *config.Config

	db storage.DB

	useCase *usecase.UseCase

	repository repository.Repository

	//grpcAuthServer usersProtos.UsersAuthServer
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
		p.cfg = config.NewConfig()
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
	}
	return p.db
}

//func (p *Provider) GRPCAuthServer() usersProtos.UsersAuthServer {
//	if p.grpcAuthServer == nil {
//		p.grpcAuthServer = grpc.NewGRPCAuthServer(p.Logger(), p.Config().GRPCServerConfig, p.UseCase())
//	}
//	return p.grpcAuthServer
//}
