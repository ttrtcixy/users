package grpcauthservise

import (
	"context"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/delivery/grpc/ports"
	"github.com/ttrtcixy/users/internal/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

type SignoutService struct {
	log     logger.Logger
	usecase ports.SignoutUseCase
}

func NewSignout(log logger.Logger, usecase ports.SignoutUseCase) *SignoutService {
	return &SignoutService{
		log:     log,
		usecase: usecase,
	}
}

func (s *SignoutService) Signout(ctx context.Context, payload *dtos.SignoutRequest) (*emptypb.Empty, error) {
	return nil, nil
}
