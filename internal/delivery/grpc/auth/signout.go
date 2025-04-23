package grpcauthservise

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthSignout interface {
	Signout(context.Context, *dtos.SignoutRequest) (*emptypb.Empty, error)
}

type signout struct {
	log     logger.Logger
	usecase *usecase.UseCase
}

func NewSignout(log logger.Logger, usecase *usecase.UseCase) AuthSignout {
	return &signout{
		log:     log,
		usecase: usecase,
	}
}

func (s *signout) Signout(ctx context.Context, payload *dtos.SignoutRequest) (*emptypb.Empty, error) {
	return nil, nil
}
