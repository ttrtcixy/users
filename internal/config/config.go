package config

import (
	"errors"
	"github.com/goloop/env"
	"github.com/ttrtcixy/users/internal/logger"
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

// NewConfig load parameters from the env file and return Config
func NewConfig(log logger.Logger) *Config {
	MustLoad(log, ".env")

	return &Config{
		DBConfig:         NewDbConfig(log),
		GRPCServerConfig: NewGRPCConfig(log),
	}
}

// MustLoad loading parameters from the env file
func MustLoad(log logger.Logger, filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatal("file %s does not exist", filename)
		} else {
			log.Fatal(err.Error())
		}
	}

	err = env.Load(filename)
	if err != nil {
		log.Fatal("Incorrect data in the configuration file: %s", err.Error())
	}

	log.Info("the configuration file '%s' has been uploaded successfully", filename)
}
