package grpcauthservise

import (
	"context"
	"errors"
	dtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/core/entities"
	"github.com/ttrtcixy/users/internal/delivery/grpc/ports"
	apperrors "github.com/ttrtcixy/users/internal/errors"
	"github.com/ttrtcixy/users/internal/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type VerifyService struct {
	log     logger.Logger
	usecase ports.VerifyUseCase
}

func NewVerify(log logger.Logger, usecase ports.UseCase) *VerifyService {
	return &VerifyService{
		log:     log,
		usecase: usecase,
	}
}

func (s *VerifyService) Verify(ctx context.Context, payload *dtos.VerifyRequest) (*dtos.VerifyResponse, error) {
	if err := s.validate(payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := s.usecase.Verify(ctx, s.DTOToEntity(payload))
	if err != nil {
		return nil, s.errResponse(err)
	}

	return s.EntityToDTO(result), nil
}

func (s *VerifyService) DTOToEntity(payload *dtos.VerifyRequest) *entities.VerifyRequest {
	return &entities.VerifyRequest{JwtToken: payload.JwtToken}
}

func (s *VerifyService) EntityToDTO(result *entities.VerifyResponse) *dtos.VerifyResponse {
	return &dtos.VerifyResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
	}
}

func (s *VerifyService) errResponse(err error) error {
	switch {
	case errors.Is(err, apperrors.ErrEmailTokenExpired):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (s *VerifyService) validate(payload *dtos.VerifyRequest) error {
	jwt := payload.GetJwtToken()

	if len(jwt) <= 0 {
		return errors.New("token required")
	}
	return nil
}
