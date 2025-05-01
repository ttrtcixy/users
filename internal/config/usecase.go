package config

import (
	"fmt"
	"os"
	"strconv"
)

type UsecaseConfig interface {
	PasswordSaltLength() int
	JWTSecret() string
}

type usecaseConfig struct {
	passwordSaltLength int
	jwtSecret          string
}

func (cfg *usecaseConfig) PasswordSaltLength() int {
	return cfg.passwordSaltLength
}

func (cfg *usecaseConfig) JWTSecret() string {
	return cfg.jwtSecret
}

func NewUsecaseConfig() (UsecaseConfig, error) {
	const op = "config.NewUsecaseConfig"

	var cfg = &usecaseConfig{}
	if env, ok := os.LookupEnv("JWT_SECRET"); ok {
		cfg.jwtSecret = env
	} else {
		//return nil, fmt.Errorf("op: %s, env variable 'JWT_SECRET' is not set", op)
	}

	if env, ok := os.LookupEnv("PASSWORD_SALT_LENGTH"); ok {
		psl, err := strconv.Atoi(env)
		if err != nil {
			return nil, fmt.Errorf("op: %s, env variable 'PASSWORD_SALT_LENGTH' bad format", op)
		}
		cfg.passwordSaltLength = psl
	} else {
		return nil, fmt.Errorf("op: %s, env variable 'PASSWORD_SALT_LENGTH' is not set", op)
	}

	return cfg, nil
}
