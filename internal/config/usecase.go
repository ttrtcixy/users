package config

import (
	"fmt"
	"os"
	"strconv"
)

//type UsecaseConfig interface {
//	PasswordSaltLength() int
//	JWTSecret() string
//}

type UsecaseConfig struct {
	passwordSaltLength int
	jwtSecret          string
}

func (cfg *UsecaseConfig) PasswordSaltLength() int {
	return cfg.passwordSaltLength
}

func (cfg *UsecaseConfig) JWTSecret() string {
	return cfg.jwtSecret
}

func (c *Config) LoadUsecaseConfig(fErr *ErrEnvVariableNotFound) {
	const op = "config.NewUsecaseConfig"

	var cfg = &UsecaseConfig{}
	if env, ok := os.LookupEnv("JWT_SECRET"); ok {
		cfg.jwtSecret = env
	} else {
		//fErr.Add(fmt.Errorf("op: %s, env variable 'JWT_SECRET' is not set", op))
	}

	if env, ok := os.LookupEnv("PASSWORD_SALT_LENGTH"); ok {
		psl, err := strconv.Atoi(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'PASSWORD_SALT_LENGTH' bad format", op))
		} else {
			cfg.passwordSaltLength = psl
		}
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'PASSWORD_SALT_LENGTH' is not set", op))
	}

	c.UsecaseConfig = cfg
}
