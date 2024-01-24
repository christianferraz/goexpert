package limiter

import (
	"context"
	"strconv"
	"time"

	"github.com/christianferraz/goexpert/Rate_Limiter/configs"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	RedisClient *redis.Client
	config      *configs.Config
}

func NewRateLimiter(config *configs.Config, client *redis.Client) *RateLimiter {
	return &RateLimiter{
		RedisClient: client,
		config:      config,
	}
}

func (r *RateLimiter) IsLimited(ctx context.Context, key string, token string) bool {
	var requestLimit int
	if isValidToken(token, r.config.AllowedTokens) {
		requestLimit = 100
	} else {
		requestLimit = 10
	}

	val, err := r.RedisClient.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return true
	}

	count, _ := strconv.Atoi(val)
	if count >= requestLimit {
		return true
	}

	count++
	if err := r.RedisClient.Set(ctx, key, strconv.Itoa(count), time.Second*time.Duration(requestLimit)).Err(); err != nil {
		return true
	}

	return false
}

func isValidToken(token string, tokenList []string) bool {
	for _, t := range tokenList {
		if t == token {
			return true
		}
	}
	return false
}
