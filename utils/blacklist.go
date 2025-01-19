package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

const blacklistedValue = "blacklisted"

func BlacklistToken(token string, redisClient *redis.Client, expiry time.Duration) error {
	ctx := context.Background()

	if err := redisClient.Set(ctx, token, blacklistedValue, expiry).Err(); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

func IsTokenBlacklisted(token string, redisClient *redis.Client) (bool, error) {
	ctx := context.Background()

	val, err := redisClient.Get(ctx, token).Result()

	// Token is not blacklisted
	if err == redis.Nil {
		return false, nil
	}

	// Redis query error
	if err != nil {
		return false, fmt.Errorf("failed to check blacklist status: %w", err)
	}

	return val == blacklistedValue, nil
}
