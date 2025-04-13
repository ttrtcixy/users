package config

import (
	"fmt"
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

func NewGRPCConfig() GRPCServerConfig {
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
