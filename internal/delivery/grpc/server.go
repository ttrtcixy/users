package grpc

import (
	"context"
	usersProtos "github.com/ttrtcixy/users-protos/gen/go/users"
	"github.com/ttrtcixy/users/internal/config"
	"github.com/ttrtcixy/users/internal/delivery/grpc/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/ttrtcixy/users/internal/logger"
	"github.com/ttrtcixy/users/internal/usecase"
)

type Server struct {
	log logger.Logger
	cfg config.GRPCServerConfig

	srv             *grpc.Server
	l               net.Listener
	userAuthService usersProtos.UsersAuthServer
	userService     usersProtos.UsersServer
}

func (s *Server) register(gRPC *grpc.Server) {
	usersProtos.RegisterUsersAuthServer(gRPC, s.userAuthService)
	usersProtos.RegisterUsersServer(gRPC, s.userService)
}

func NewGRPCServer(log logger.Logger, cfg config.GRPCServerConfig, usecase *usecase.UseCase) *Server {
	return &Server{
		log:             log,
		cfg:             cfg,
		userAuthService: grpcauthservise.NewUserAuthService(context.Background(), log, usecase),
	}
}

func (s *Server) Start(ctx context.Context) (err error) {
	s.log.Info("[*] starting gRPC server on %s", s.cfg.Addr())

	s.srv = grpc.NewServer()
	s.register(s.srv)
	s.l, err = net.Listen(s.cfg.Network(), s.cfg.Addr())
	if err != nil {
		return err
	}
	reflection.Register(s.srv)
	return s.srv.Serve(s.l)
}

func (s *Server) Close() error {
	s.srv.Stop()

	return nil
}
