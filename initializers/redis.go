package initializers

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var RedisClient *redis.Client

func ConnectRedis(config *Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPass,
	})

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatal("Could not connect to Redis")
	}

	fmt.Println("✔ Successfully connected to Redis client.")
}
