package config

import (
	"fmt"
	"os"
)

type DSNConfig struct {
	Username string
	Password string
	Database string
	Host     string
	Port     string
}

type DBConfig struct {
	Driver string
	DSN    DSNConfig
}

func (cfg DSNConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func NewDatabaseConfig() *DSNConfig {
	return &DSNConfig{
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	}
}
