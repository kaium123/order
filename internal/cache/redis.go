package cache

import (
	"github.com/go-redis/redis/v8"
)

type Config struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func New(config *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	return rdb
}
