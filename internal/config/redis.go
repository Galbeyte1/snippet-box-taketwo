package config

import (
	"os"
	"time"
)

type RedisConfig struct {
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

func NewRedisConfig() *RedisConfig {
	// Parse Redis timeouts
	idleTimeoutSeconds := parseEnvInt("REDIS_IDLE_TIMEOUT", 240)
	redisIdleTimeout := time.Duration(idleTimeoutSeconds) * time.Second

	return &RedisConfig{
		Addr:        os.Getenv("REDIS_ADDR"),
		MaxIdle:     parseEnvInt("REDIS_MAX_IDLE", 10),
		MaxActive:   parseEnvInt("REDIS_MAX_ACTIVE", 100),
		IdleTimeout: redisIdleTimeout,
	}

}
