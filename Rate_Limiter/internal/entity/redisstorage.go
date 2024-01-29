package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	ctx    context.Context
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		ctx:    context.Background(),
		client: client,
	}
}

func (r *RedisStorage) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisStorage) Update(key string, value string) error {
	// Obter o TTL atual da chave
	ttlVal, err := r.client.TTL(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error getting key %s ttl: %w", key, err)
	}

	// Atualizar o valor da chave
	var ttl time.Duration
	if ttlVal > 0 {
		ttl = ttlVal
	} else {
		ttl = 0 // Ou a duração padrão que você deseja
	}
	err = r.client.Set(r.ctx, key, value, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}

	return nil
}

func (r *RedisStorage) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}
