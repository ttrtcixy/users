package config

import (
	"github.com/ttrtcixy/users/internal/logger"
	"os"
)

type DBConfig interface {
	DSN() string
}

type dbConfig struct {
	path string
}

func NewDbConfig(log logger.Logger) DBConfig {
	var cfg = &dbConfig{}

	if value, ok := os.LookupEnv("DB_PATH"); ok {
		cfg.path = value
	} else {
		log.Error("[!] DB_PATH is not set")
	}

	return cfg
}

func (c *dbConfig) DSN() string {
	return c.path
}
