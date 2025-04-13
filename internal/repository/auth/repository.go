package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/storage"
)

type AuthRepository interface {
	CheckUserPassword(ctx context.Context, user *entities.User) (*entities.User, error)
}

type authRepository struct {
	log logger.Logger
	DB  storage.DB
}

func NewAuthRepository(ctx context.Context, log logger.Logger, db storage.DB) AuthRepository {
	return &authRepository{
		log: log,
		DB:  db,
	}
}
