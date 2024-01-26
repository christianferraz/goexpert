package limiter

import (
	"context"

	"github.com/christianferraz/goexpert/Rate_Limiter/configs"
	"github.com/go-redis/redis_rate/v10"
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
	r_limiter := redis_rate.NewLimiter(r.RedisClient)
	var requestLimit int
	if isValidToken(token, r.config.AllowedTokens) {
		requestLimit = 100
	} else {
		requestLimit = 10
	}

	requestLimit++
	res, err := r_limiter.Allow(ctx, key, redis_rate.PerSecond(requestLimit))
	if err != nil {
		return true
	}
	if res.Remaining <= 0 {
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
