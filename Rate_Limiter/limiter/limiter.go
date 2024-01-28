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

func (r *RateLimiter) IsLimited(ctx context.Context, ip string, key string, token string) bool {
	r_limiter := redis_rate.NewLimiter(r.RedisClient)

	requestLimit := isValidToken(ip, token, r.config.AllowedTokens)

	res, err := r_limiter.Allow(ctx, key, redis_rate.PerSecond(requestLimit))
	if err != nil {
		return true
	}
	if res.Remaining <= 0 {
		return true
	}
	return false
}

func isValidToken(token string, ip string, tokenList map[string]int) int {
	for k, v := range tokenList {
		if k == token || k == ip {
			return v
		}
	}
	return 10
}
