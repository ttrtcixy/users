package repository

import (
	"context"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/repository/auth"
	"github.com/ttrtcixy/users/internal/storage"
)

type Repository struct {
	*authrepo.AuthRepository
}

func NewRepository(ctx context.Context, log logger.Logger, db storage.DB) *Repository {
	return &Repository{
		authrepo.NewAuthRepository(ctx, log, db),
	}
}
