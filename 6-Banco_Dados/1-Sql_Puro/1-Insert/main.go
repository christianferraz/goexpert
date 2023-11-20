package main

import (
	"database/sql"

	"github.com/christianferraz/goexpert/6-Banco_Dados/1-Sql_Puro/1-Insert/internal/entity"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}

	defer db.Close()
	product := entity.NewProduct(db)
	product, err = product.InsertProduct("Product 12", 10.0)
	if err != nil {
		panic(err)
	}
	println(product.ID, product.Name, product.Price)
}
