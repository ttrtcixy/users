package authusecase

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Signout interface {
	Run(ctx context.Context, payload *dtos.SignoutRequest) (*emptypb.Empty, error)
}

type signoutRepository interface {
}

type signout struct {
	log  logger.Logger
	repo signoutRepository
}

func NewSignout(ctx context.Context, log logger.Logger, repo signoutRepository) Signout {
	return &signout{
		log:  log,
		repo: repo,
	}
}

func (l *signout) Run(ctx context.Context, payload *dtos.SignoutRequest) (*emptypb.Empty, error) {
	return nil, nil
}
