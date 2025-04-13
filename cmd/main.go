package main

import (
	"context"
	"github.com/ttrtcixy/users/internal/app"
)

func main() {
	a := app.NewApp()

	a.Run(context.Background())
}
