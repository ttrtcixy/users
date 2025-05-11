package repository

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/repository/auth"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/storage/pg"
)

type Repository struct {
	*authrepo.AuthRepository
}

func NewRepository(ctx context.Context, log logger.Logger, db storage.DB) *Repository {
	return &Repository{
		authrepo.NewAuthRepository(ctx, log, db),
	}
}

func (r *Repository) RunInTx(ctx context.Context, fn func(context.Context) error) error {
	err := r.DB.RunInTx(ctx, fn)
	if err != nil {
		return err
	}
	return nil
}
