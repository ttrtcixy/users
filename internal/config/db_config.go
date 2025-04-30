package config

import (
	"errors"
	"os"
)

type DBConfig interface {
	DSN() string
}

type dbConfig struct {
	dsn string
}

func NewDbConfig() (DBConfig, error) {
	const op = "config.NewDbConfig"
	var cfg = &dbConfig{}

	if value, ok := os.LookupEnv("DB_URL"); ok {
		cfg.dsn = value
	} else {
		return nil, errors.New("op: " + op + ", env parameter 'DB_URL' is not set")
	}

	return cfg, nil
}

func (c *dbConfig) DSN() string {
	return c.dsn
}
