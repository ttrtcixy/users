package config

import (
	"errors"
	"fmt"
	"github.com/goloop/env"
	"os"
)

// Config struct
type Config struct {
	DBConfig         DBConfig
	GRPCServerConfig GRPCServerConfig
}

func (c *Config) Close() error {
	env.Clear()
	return nil
}

// New load parameters from the env file and return Config
func New() (*Config, error) {
	err := MustLoad(".env")
	if err != nil {
		return nil, err
	}

	dbCfg, err := NewDbConfig()
	if err != nil {
		return nil, err
	}

	grpcGfg, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		DBConfig:         dbCfg,
		GRPCServerConfig: grpcGfg,
	}, nil
}

// MustLoad loading parameters from the env file
func MustLoad(filename string) error {
	const op = "config.MustLoad"
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("op: %s, file: %s does not exist", op, filename)
		} else {
			return err
		}
	}

	err = env.Load(filename)
	if err != nil {
		return fmt.Errorf("op: %s,incorrect data in the configuration file: %s", op, err.Error())
	}
	return nil
}
