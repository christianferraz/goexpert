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
	requestLimit := isValidToken(key, r.config.AllowedTokens)

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
		return true // Requisição limitada
	}
	currentCount++
	// Incrementar contagem
	fmt.Println("err currentCount:", currentCount)
	err = r.strClient.Update(key, strconv.Itoa(currentCount))
	return err != nil
}

func isValidToken(key string, tokenList map[string]int) int {
	if limit, exists := tokenList[key]; exists {
		return limit
	}
	return 10
}
