package middleware

import (
	"context"
	"fmt"
	"github.com/ttrtcixy/users/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

func RecoveryUnaryInterceptor(log logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func(log logger.Logger) {
			if r := recover(); r != nil {
				e := fmt.Errorf("panic in %s: %v", info.FullMethod, r)
				log.Error(e.Error())
				debug.PrintStack()
				err = status.Error(codes.Internal, "internal server error")
			}
		}(log)
		return handler(ctx, req)
	}
}
