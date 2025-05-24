package database

import (
	"github.com/Galbeyte1/snippet-box-taketwo/internal/config"
	"github.com/gomodule/redigo/redis"
)

func OpenRedis(cfg config.RedisConfig) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Addr)
		},
	}
}
