package config

import (
	"fmt"
	"os"
)

//type GRPCServerConfig interface {
//	Addr() string
//	Port() string
//	Host() string
//	Network() string
//}

type GRPCServerConfig struct {
	host    string
	port    string
	network string
}

func (c *Config) LoadGRPCConfig(fErr *ErrEnvVariableNotFound) {
	const op = "config.NewGRPCConfig"
	cfg := &GRPCServerConfig{}

	if value, ok := os.LookupEnv("GRPC_HOST"); ok {
		cfg.host = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_HOST' is not set", op))
	}

	if value, ok := os.LookupEnv("GRPC_PORT"); ok {
		cfg.port = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_PORT' is not set", op))
	}

	if value, ok := os.LookupEnv("GRPC_NETWORK"); ok {
		cfg.network = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'GRPC_NETWORK' is not set", op))
	}

	c.GRPCServerConfig = cfg
}

func (c *GRPCServerConfig) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func (c *GRPCServerConfig) Port() string {
	return c.port
}

func (c *GRPCServerConfig) Host() string {
	return c.host
}

func (c *GRPCServerConfig) Network() string {
	return c.network
}
