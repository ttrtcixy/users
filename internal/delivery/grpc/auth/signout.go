package grpcauthservise

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/delivery/grpc/ports"
	"github.com/ttrtcixy/users/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if err := s.validate(payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.usecase.Signout(ctx, s.DTOToEntity(payload)); err != nil {
		return nil, s.errResponse(err)
	}

	return nil, nil
}

func (s *SignoutService) DTOToEntity(payload *dtos.SignoutRequest) *entities.SignoutRequest {
	return &entities.SignoutRequest{ClientUUID: payload.ClientId}
}

func (s *SignoutService) errResponse(err error) error {
	return status.Error(codes.Internal, err.Error())
}

func (s *SignoutService) validate(payload *dtos.SignoutRequest) error {
	if err := uuid.Validate(payload.GetClientId()); err != nil {
		return fmt.Errorf("invalid client_id")
	}

	return nil
}
