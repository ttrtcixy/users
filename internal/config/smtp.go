package config

import (
	"fmt"
	"os"
)

type SmtpConfig struct {
	host     string
	port     string
	secure   string
	sender   string
	password string
}

func (c *Config) LoadSmtpConfig(fErr *ErrEnvVariableNotFound) {
	const op = "Config.LoadSmtpConfig"
	cfg := &SmtpConfig{}

	if value, ok := os.LookupEnv("SMTP_HOST"); ok {
		cfg.host = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_HOST' is not set", op))
	}

	if value, ok := os.LookupEnv("SMTP_PORT"); ok {
		cfg.port = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_PORT' is not set", op))
	}

	if value, ok := os.LookupEnv("SMTP_SECURE"); ok {
		cfg.secure = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_SECURE' is not set", op))
	}

	if value, ok := os.LookupEnv("SMTP_SENDER"); ok {
		cfg.sender = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_SENDER' is not set", op))
	}

	if value, ok := os.LookupEnv("SMTP_PASSWORD"); ok {
		cfg.password = value
	} else {
		fErr.Add(fmt.Errorf("%s: env variable 'SMTP_PASSWORD' is not set", op))
	}

	c.SmtpConfig = cfg
}
func (c *SmtpConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host(), c.Port())
}

func (c *SmtpConfig) Host() string {
	return c.host
}

func (c *SmtpConfig) Port() string {
	return c.port
}

func (c *SmtpConfig) Password() string {
	return c.password
}

func (c *SmtpConfig) Sender() string {
	return c.sender
}
