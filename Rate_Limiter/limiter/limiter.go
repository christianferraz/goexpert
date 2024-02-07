package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/christianferraz/goexpert/Rate_Limiter/configs"
	"github.com/christianferraz/goexpert/Rate_Limiter/internal/entity"
)

type RateLimiter struct {
	strClient entity.Storage
	config    *configs.Config
}

func NewRateLimiter(config *configs.Config, str entity.Storage) *RateLimiter {
	return &RateLimiter{
		strClient: str,
		config:    config,
	}
}

func (r *RateLimiter) IsLimited(ctx context.Context, key string) (bool, error) {
	requestLimit, time_block := isValidToken(key, r.config.AllowedTokens)
	if time_block > 0 {
		value, err := r.strClient.Get(fmt.Sprintf("%s-%s", "block", key))
		if err != nil {
			return false, err
		}
		if value != "" {
			return true, nil
		}
	}

	currentCountStr, err := r.strClient.Get(key)
	if err != nil {
		err = r.strClient.Set(key, strconv.Itoa(1), time.Second*1)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	currentCount, err := strconv.Atoi(currentCountStr)
	if err != nil {
		return false, err
	}
	if currentCount > requestLimit {
		if time_block > 0 {
			key := fmt.Sprintf("%s-%s", "block", key)
			err = r.strClient.Set(key, strconv.Itoa(time_block), time.Second*time.Duration(time_block))
			if err != nil {
				return false, err
			}
		}
		return true, nil
	}
	currentCount++
	err = r.strClient.Update(key, strconv.Itoa(currentCount))
	if err != nil {
		return true, err
	}
	return false, nil
}

func isValidToken(key string, tokenList map[string][]int) (int, int) {
	if limit, exists := tokenList[key]; exists {
		return limit[0], limit[1]
	}
	return 10, 0
}
