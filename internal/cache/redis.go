package cache

import (
	"context"
	"log"
	"time"

	"github.com/arslan-atajykov/shortener-api/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init(cfg *config.Config) {
	Client = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis")

}
