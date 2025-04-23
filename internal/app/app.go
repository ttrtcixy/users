package app

import (
	"context"
	"sync"
)

type App struct {
	wg sync.WaitGroup
	*Provider
}

func NewApp() *App {
	return &App{
		Provider: NewProvider(),
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
