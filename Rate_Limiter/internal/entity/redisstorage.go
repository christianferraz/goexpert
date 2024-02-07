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
	// Obter o TTL atual da chave em milissegundos
	ttlVal, err := r.client.PTTL(r.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error getting TTL for key %s: %w", key, err)
	}
	// Atualizar o valor da chave
	err = r.client.Set(r.ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("error setting key %s: %w", key, err)
	}

	// Se a chave tinha um TTL, restaurar o TTL
	if ttlVal > 0 {
		err = r.client.PExpire(r.ctx, key, ttlVal).Err() // Usar PExpire para manter precisão em milissegundos
		if err != nil {
			return fmt.Errorf("error setting TTL for key %s: %w", key, err)
		}
	} else {
		r.client.Del(r.ctx, key)
		return fmt.Errorf("key %s has no TTL", key)
	}

	return nil
}

func (r *RedisStorage) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}
