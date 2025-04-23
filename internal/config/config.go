package config

import (
	"errors"
	"github.com/goloop/env"
	"github.com/ttrtcixy/users/internal/logger"
	"log"
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
	MustLoad(".env")

	return &Config{
		DBConfig:         NewDbConfig(log),
		GRPCServerConfig: NewGRPCConfig(log),
	}
}

// MustLoad loading parameters from the env file
func MustLoad(filename string) {
	_, err := os.Stat(filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("file %s does not exist", filename)
		} else {
			log.Fatal(err)
		}
	}

	err = env.Load(filename)
	if err != nil {
		log.Fatalf("Incorrect data in the configuration file: %s", err)
	}

	log.Printf("the configuration file '%s' has been uploaded successfully", filename)
}
