package config

import (
	"fmt"
	"github.com/ttrtcixy/users/internal/logger"
	"os"
)

type GRPCServerConfig interface {
	Addr() string
	Port() string
	Host() string
	Network() string
}

type gRPCServerConfig struct {
	host    string
	port    string
	network string
}

func NewGRPCConfig(log logger.Logger) GRPCServerConfig {
	var cfg = &gRPCServerConfig{}

	if value, ok := os.LookupEnv("GRPC_HOST"); ok {
		cfg.host = value
	} else {
		log.Error("[!] GRPC_HOST is not set")
	}

	if value, ok := os.LookupEnv("GRPC_PORT"); ok {
		cfg.port = value
	} else {
		log.Error("[!] GRPC_PORT is not set")
	}

	if value, ok := os.LookupEnv("GRPC_NETWORK"); ok {
		cfg.network = value
	} else {
		log.Error("[!] GRPC_NETWORK is not set")
	}

	return &gRPCServerConfig{host: os.Getenv("GRPC_HOST"), port: os.Getenv("GRPC_PORT"), network: os.Getenv("GRPC_NETWORK")}
}

func (c *gRPCServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *gRPCServerConfig) Port() string {
	return c.port
}

func (c *gRPCServerConfig) Host() string {
	return c.host
}

func (c *gRPCServerConfig) Network() string {
	return c.network
}
