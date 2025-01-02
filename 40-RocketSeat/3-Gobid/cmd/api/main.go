package main

import (
	"context"

	"github.com/christianferraz/goexpert/40-RocketSeat/3-Gobid/internal/api"
)

func main() {
	ctx := context.Background()
	api := api.NewApi(ctx)
	api.BindRoutes()
	api.Start()
}
