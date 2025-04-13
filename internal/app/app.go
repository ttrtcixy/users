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
	gRPCServer grpc.GRPCServer
}

func NewApp() *App {
	return &App{
		provider: NewProvider(),
	}
}

func (a *App) Run(ctx context.Context) {
	defer globalCloser.CloseAll()

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
	return a.gRPCServer.Start(ctx, a.provider.Config().GRPCServerConfig)
}
