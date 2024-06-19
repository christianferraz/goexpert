package main

import (
	"context"
	"log"

	"github.com/christianferraz/goexpert/26-Leilao/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		panic(err)
	}

	databaseClient, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	defer databaseClient.Client().Disconnect(ctx)
}
