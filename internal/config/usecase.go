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

	smtpHost   string
	smtpPort   int
	smtpSecure string
	smtpSender string
	password   string
}

func (cfg *UsecaseConfig) PasswordSaltLength() int {
	return cfg.passwordSaltLength
}

func (cfg *UsecaseConfig) JWTSecret() string {
	return cfg.jwtSecret
}

func (cfg *UsecaseConfig) SMTPAddr() string {
	return fmt.Sprintf("%s:%d", cfg.SMTPHost(), cfg.SMTPPort())
}

func (cfg *UsecaseConfig) SMTPHost() string {
	return cfg.smtpHost
}

func (cfg *UsecaseConfig) SMTPPort() int {
	return cfg.smtpPort
}

func (cfg *UsecaseConfig) SMTPSecure() string {
	return cfg.smtpSecure
}

func (cfg *UsecaseConfig) SMTPSender() string {
	return cfg.smtpSender
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

	cfg.SMTPLoad(fErr)

	c.UsecaseConfig = cfg
}

func (cfg *UsecaseConfig) SMTPLoad(fErr *ErrEnvVariableNotFound) {
	const op = "config.NewUsecaseConfig.SMTP"
	if env, ok := os.LookupEnv("SMTP_HOST"); ok {
		cfg.smtpHost = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_HOST' is not set", op))
	}

	if env, ok := os.LookupEnv("SMTP_SECURE"); ok {
		cfg.smtpSecure = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_SECURE' is not set", op))
	}

	if env, ok := os.LookupEnv("SMTP_PORT"); ok {
		psl, err := strconv.Atoi(env)
		if err != nil {
			fErr.Add(fmt.Errorf("%s: env variable 'SMTP_PORT' bad format", op))
		} else {
			cfg.smtpPort = psl
		}
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_PORT' is not set", op))
	}

	if env, ok := os.LookupEnv("SMTP_SENDER"); ok {
		cfg.smtpSender = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_SENDER' is not set", op))
	}

	if env, ok := os.LookupEnv("SMTP_PASSWORD"); ok {
		cfg.password = env
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_PASSWORD' is not set", op))
	}
}

func (cfg *UsecaseConfig) SMTPLogin() string {
	return cfg.smtpSender
}

func (cfg *UsecaseConfig) SMTPPassword() string {
	return cfg.password
}
