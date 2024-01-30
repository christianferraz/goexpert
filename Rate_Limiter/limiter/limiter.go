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

func (r *RateLimiter) IsLimited(ctx context.Context, key string) bool {
	requestLimit, time_block := isValidToken(key, r.config.AllowedTokens)
	if time_block > 0 {
		value, _ := r.strClient.Get(fmt.Sprintf("%s-%s", "block", key))
		if value != "" {
			return true
		}
	}

	currentCountStr, err := r.strClient.Get(key)

	if err != nil {
		err = r.strClient.Set(key, strconv.Itoa(1), time.Second*1)
		return err != nil
	}

	currentCount, err := strconv.Atoi(currentCountStr)
	if err != nil {
		// Tratar erro de conversão
		return true
	}
	if currentCount >= requestLimit {
		if time_block > 0 {
			key := fmt.Sprintf("%s-%s", "block", key)
			r.strClient.Set(key, strconv.Itoa(time_block), time.Second*time.Duration(time_block))
		}
		return true // Requisição limitada
	}
	currentCount++
	// Incrementar contagem
	fmt.Println("err currentCount:", currentCount)
	err = r.strClient.Update(key, strconv.Itoa(currentCount))
	return err != nil
}

func isValidToken(key string, tokenList map[string][]int) (int, int) {
	if limit, exists := tokenList[key]; exists {
		return limit[0], limit[1]
	}
	return 10, 0
}
