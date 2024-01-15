package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/christianferraz/goexpert/Rate_Limiter/configs"
	"github.com/christianferraz/goexpert/Rate_Limiter/middleware"
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
	defer r.Close()
	pong, err := r.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Conexão ao Redis estabelecida:", pong)

	http.HandleFunc("/", middleware.CountMiddleware(handler, &ctx, r))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}
