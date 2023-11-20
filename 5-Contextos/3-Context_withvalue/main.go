package main

import "context"

type keyString string

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyString("userID"), 123)
	bookHotel(ctx)
}

func bookHotel(ctx context.Context) {
	userID := ctx.Value(keyString("userID")).(int)
	println(userID)
}
