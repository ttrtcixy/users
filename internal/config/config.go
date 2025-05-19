package config

import (
	"errors"
	"fmt"
	"github.com/goloop/env"
	"os"
	"strings"
)

type ErrEnvVariableNotFound struct {
	Variables []error
}

func (e *ErrEnvVariableNotFound) Error() string {
	if len(e.Variables) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("missing or invalid configuration:\n")
	for _, err := range e.Variables {
		sb.WriteString(" - ")
		sb.WriteString(err.Error())
		sb.WriteString("\n")
	}
	return sb.String()
}

func (e *ErrEnvVariableNotFound) Add(err error) {
	e.Variables = append(e.Variables, err)
}

// Config struct
type Config struct {
	DBConfig         *DBConfig
	GRPCServerConfig *GRPCServerConfig
	UsecaseConfig    *UsecaseConfig
	SmtpConfig       *SmtpConfig
	CloserConfig     *CloserConfig
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
	var cfg = &Config{}

	var envErrs = &ErrEnvVariableNotFound{}

	cfg.LoadDbConfig(envErrs)

	cfg.LoadGRPCConfig(envErrs)

	cfg.LoadUsecaseConfig(envErrs)

	cfg.LoadSmtpConfig(envErrs)

	cfg.LoadCloserConfig(envErrs)

	if envErrs.Variables != nil {
		return nil, envErrs
	}

	return cfg, nil
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
