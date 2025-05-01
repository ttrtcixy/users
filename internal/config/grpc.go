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

func NewGRPCConfig() (GRPCServerConfig, error) {
	const op = "config.NewGRPCConfig"
	var cfg = &gRPCServerConfig{}

	if value, ok := os.LookupEnv("GRPC_HOST"); ok {
		cfg.host = value
	} else {
		return nil, fmt.Errorf("op: %s, env variable 'GRPC_HOST' is not set", op)
	}

	if value, ok := os.LookupEnv("GRPC_PORT"); ok {
		cfg.port = value
	} else {
		return nil, fmt.Errorf("op: %s, env variable 'GRPC_PORT' is not set", op)
	}

	if value, ok := os.LookupEnv("GRPC_NETWORK"); ok {
		cfg.network = value
	} else {
		return nil, fmt.Errorf("op: %s, env variable 'GRPC_NETWORK' is not set", op)
	}

	return cfg, nil
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
