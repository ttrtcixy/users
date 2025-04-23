package app

import (
	"context"
	"github.com/ttrtcixy/users/internal/delivery/grpc"
	"log"
	"sync"
)

type App struct {
	wg         sync.WaitGroup
	provider   *Provider
	gRPCServer *grpc.Server
}

func NewApp() *App {
	return &App{
		provider: NewProvider(),
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.provider.Closer().Close()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		err := a.startGRPCServer(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	a.wg.Wait()
}

func (a *App) startGRPCServer(ctx context.Context) error {
	a.gRPCServer = a.provider.GRPCServer()

	return a.gRPCServer.Start(ctx, a.provider.Config().GRPCServerConfig)
}
