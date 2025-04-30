package authrepo

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
)

func (r *AuthRepository) CheckUserPassword(ctx context.Context, user *entities.User) (bool, error) {
	return false, nil
}
