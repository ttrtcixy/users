package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

//type UsecaseConfig interface {
//	PasswordSaltLength() int
//	JWTSecret() string
//}

type UsecaseConfig struct {
	passwordSaltLength int
	jwtSecret          string
	emailJwtExpiry     time.Duration
	accessJwtExpiry    time.Duration
	refreshJwtExpiry   time.Duration
}

func (cfg *UsecaseConfig) PasswordSaltLength() int {
	return cfg.passwordSaltLength
}

func (cfg *UsecaseConfig) JWTSecret() string {
	return cfg.jwtSecret
}

func (cfg *UsecaseConfig) EmailJwtExpiry() time.Duration {
	return cfg.emailJwtExpiry
}

func (cfg *UsecaseConfig) AccessJwtExpiry() time.Duration {
	return cfg.accessJwtExpiry
}

func (cfg *UsecaseConfig) RefreshJwtExpiry() time.Duration {
	return cfg.refreshJwtExpiry
}

func (c *Config) LoadUsecaseConfig(fErr *ErrEnvVariableNotFound) {
	const op = "LoadUsecaseConfig"

	var cfg = &UsecaseConfig{}
	if env, ok := os.LookupEnv("JWT_SECRET"); ok {
		cfg.jwtSecret = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'JWT_SECRET' is not set", op))
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

	if env, ok := os.LookupEnv("EMAIL_JWT_EXPIRY"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'EMAIL_JWT_EXPIRY' bad format", op))
		} else {
			cfg.emailJwtExpiry = t
		}
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'EMAIL_JWT_EXPIRY' is not set", op))
	}

	if env, ok := os.LookupEnv("ACCESS_JWT_EXPIRY"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'ACCESS_JWT_EXPIRY' bad format", op))
		} else {
			cfg.accessJwtExpiry = t
		}
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'ACCESS_JWT_EXPIRY' is not set", op))
	}

	if env, ok := os.LookupEnv("REFRESH_JWT_EXPIRY"); ok {
		t, err := time.ParseDuration(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'REFRESH_JWT_EXPIRY' bad format", op))
		} else {
			cfg.refreshJwtExpiry = t
		}
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'REFRESH_JWT_EXPIRY' is not set", op))
	}

	c.UsecaseConfig = cfg
}
