package grpc

import (
	"context"
	usersProtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/delivery/grpc/auth"
	"google.golang.org/grpc"
	"net"

	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type GRPCServer struct {
	log logger.Logger
	cfg config.GRPCServerConfig

	srv *grpc.Server
	l   net.Listener

	userAuthService usersProtos.UsersAuthServer
}

func (s *GRPCServer) register(gRPC *grpc.Server) {
	usersProtos.RegisterUsersAuthServer(gRPC, s.userAuthService)
}

func NewGRPCServer(log logger.Logger, cfg config.GRPCServerConfig, usecase *usecase.UseCase) *GRPCServer {
	return &GRPCServer{
		log:             log,
		cfg:             cfg,
		userAuthService: grpcauthservise.NewUserAuthService(context.Background(), log, usecase),
	}
}

func (s *GRPCServer) Start(ctx context.Context, cfg config.GRPCServerConfig) (err error) {
	s.srv = grpc.NewServer()
	s.register(s.srv)
	s.l, err = net.Listen(cfg.Network(), cfg.Addr())
	if err != nil {
		return err
	}
	return s.srv.Serve(s.l)
}

func (s *GRPCServer) Close() error {
	s.log.Info("stopping grpc server")
	s.srv.Stop()
	s.log.Info("close net listener")
	err := s.l.Close()
	if err != nil {
		return err
	}
	return nil
}
