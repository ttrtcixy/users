package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/core/entities"
)

func (r *AuthRepository) CheckUserPassword(ctx context.Context, user *entities.User) (bool, error) {
	return false, nil
}
