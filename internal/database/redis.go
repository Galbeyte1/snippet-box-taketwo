package database

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

type RedisConfig struct {
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

func (cfg RedisConfig) OpenRedis() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.Addr)
		},
	}
}
