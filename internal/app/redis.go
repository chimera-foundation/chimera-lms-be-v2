package app

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewRedis(v *viper.Viper, log *logrus.Logger) *redis.Client {
	addr := v.GetString("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: v.GetString("REDIS_PASSWORD"), 
		DB:       v.GetInt("REDIS_DB"),       
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Info("Redis connection established successfully")
	return rdb
}