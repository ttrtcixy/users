package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/storage"
)

type AuthRepository struct {
	log logger.Logger
	DB  storage.DB
}

func NewAuthRepository(ctx context.Context, log logger.Logger, db storage.DB) *AuthRepository {
	return &AuthRepository{
		log: log,
		DB:  db,
	}
}
