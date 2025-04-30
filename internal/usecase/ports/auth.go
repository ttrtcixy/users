package ports

import (
	"context"
	"github.com/ttrtcixy/users/internal/entities"
)

type Repository interface {
	AuthRepository
}

type AuthRepository interface {
	SignupRepository
	//SigninRepository
}

type SignupRepository interface {
	CheckLoginExist(ctx context.Context, payload *entities.SignupRequest) (*entities.CheckLoginResponse, error)
}

//type SigninRepository interface {
//	User(ctx context.Context, user *entities.User) (*entities.User, error)
//}
