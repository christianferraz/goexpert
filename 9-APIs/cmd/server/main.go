package main

import (
	"net/http"

	"github.com/christianferraz/goexpert/9-APIs/configs"
	"github.com/christianferraz/goexpert/9-APIs/internal/entity"
	"github.com/christianferraz/goexpert/9-APIs/internal/infra/database"
	"github.com/christianferraz/goexpert/9-APIs/internal/webserver/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.User{}, &entity.Product{})
	productDB := database.NewProductDB(db)
	productHandler := handler.NewProductHandler(productDB)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", r)
}
