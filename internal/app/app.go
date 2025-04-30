package app

import (
	"context"
	"github.com/ttrtcixy/users/internal/app/provider"
	"sync"
)

type App struct {
	wg sync.WaitGroup
	*provider.Provider
}

func NewApp() *App {
	return &App{
		Provider: provider.NewProvider(),
	}
}

func (a *App) Run(ctx context.Context) {
	defer a.Closer().Close()

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		err := a.GRPCServer().Start(ctx)
		if err != nil {
			a.Logger().Error(err.Error())
		}
	}()

	a.wg.Wait()
}
