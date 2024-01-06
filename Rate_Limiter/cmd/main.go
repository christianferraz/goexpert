package main

import (
	"context"
	"fmt"

	"github.com/christianferraz/goexpert/Rate_Limiter/configs"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := configs.LoadConfig(".")
	ctx := context.Background()
	if err != nil {
		panic(err)
	}
	r := redis.NewClient(&redis.Options{
		Addr:     config.RedisSrc,
		Password: config.RedisPass,
		DB:       0,
	})
	pong, err := r.Ping(ctx).Result()
	if err != nil {
		panic("eu" + err.Error())
	}
	fmt.Println("Conexão ao Redis estabelecida:", pong)
}
