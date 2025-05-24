package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIAddr     string
	StaticDir   string
	Database    *DSNConfig
	Redis       *RedisConfig
	SessionKey  string
	SessionOpts *SessionConfig
	Env         string
	TLSOpts     *TLSConfig
}

func LoadConfigFromEnv() Config {

	env := os.Getenv("ENV")

	// Build config
	cfg := Config{
		APIAddr:     os.Getenv("API_ADDR"),
		StaticDir:   os.Getenv("STATIC_DIR"),
		SessionKey:  os.Getenv("SESSION_SECRET"),
		Env:         env,
		Database:    NewDatabaseConfig(),
		Redis:       NewRedisConfig(),
		SessionOpts: NewSessionConfig(),
		TLSOpts:     NewTLSConfig(),
	}

	if cfg.IsProduction() {
		cfg.SessionOpts.Secure = true
	}

	return cfg
}

func parseEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return i
}

func (c Config) IsProduction() bool {
	return c.Env == "production"
}
