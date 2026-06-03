package cache

import (
	"context"
	"log"
	"time"

	"ai-nutrition/backend/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("redis unavailable, continuing without cache: %v", err)
		return nil
	}

	return client
}
