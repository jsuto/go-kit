package redisclient

import (
    "crypto/tls"
    "time"
    "github.com/redis/go-redis/v9"
)

type RedisConfig struct {
    Addr         string
    Password     string
    DB           int
    PoolSize     int
    MinIdleConns int
    DialTimeout  time.Duration
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
    PoolTimeout  time.Duration
    UseTLS       bool
}

func NewRedisClient(cfg RedisConfig) *redis.Client {
    options := &redis.Options{
        Addr:         cfg.Addr,
        Password:     cfg.Password,
        DB:           cfg.DB,
        PoolSize:     cfg.PoolSize,
        MinIdleConns: cfg.MinIdleConns,
        DialTimeout:  cfg.DialTimeout,
        ReadTimeout:  cfg.ReadTimeout,
        WriteTimeout: cfg.WriteTimeout,
        PoolTimeout:  cfg.PoolTimeout,
    }

    if cfg.UseTLS {
        options.TLSConfig = &tls.Config{
            InsecureSkipVerify: true, // or false if you validate certs
        }
    }

    return redis.NewClient(options)
}
